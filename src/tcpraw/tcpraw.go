package tcpraw

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"
	"time"

	gopacket "github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
	pcap "github.com/google/gopacket/pcap"
)

var (
	errOpNotImplemented = errors.New("operation not implemented")
	source              = rand.NewSource(time.Now().UnixNano())
	expire              = 30 * time.Minute
)

type message struct {
	bts  []byte
	addr net.Addr
}

type tcpFlow struct {
	handle       *pcap.Handle
	ready        chan struct{}
	seq          uint32
	ack          uint32
	linkLayer    gopacket.SerializableLayer
	networkLayer gopacket.SerializableLayer
	ts           time.Time
}

type TCPConn struct {
	die     chan struct{}
	dieOnce sync.Once

	tcpconn     *net.TCPConn
	listener    *net.TCPListener
	osConns     map[string]net.Conn
	osConnsLock sync.Mutex

	handles      []*pcap.Handle
	packetSource *gopacket.PacketSource
	chMessage    chan message
	flowTable    map[string]tcpFlow
	flowsLock    sync.Mutex
}

func (conn *TCPConn) lockflow(addr net.Addr, f func(e *tcpFlow)) {
	key := addr.String()
	conn.flowsLock.Lock()
	e, ok := conn.flowTable[key]
	e, ok := conn.flowTable[key]
	if !ok { // entry first visit
		e.ready = make(chan struct{})
	}
	f(&e)
	conn.flowTable[key] = e
	conn.flowsLock.Unlock()
}

// clear expired conns for listener
func (conn *TCPConn) cleaner() {
	ticker := time.NewTicker(time.Minute)
	select {
	case <-conn.die:
		return
	case <-ticker.C:
		conn.flowsLock.Lock()
		for k, v := range conn.flowTable {
			if time.Now().Sub(v.ts) > expire {
				delete(conn.flowTable, k)
			}
		}
		conn.flowsLock.Unlock()
	}
}

func (conn *TCPConn) captureFlow(handle *pcap.Handle) {
	source := gopacket.NewPacketSource(handle, handle.LinkType())

	go func() {
		for packet := range source.Packets() {
			transport := packet.TransportLayer().(*layers.TCP)

			// build address
			var ip []byte
			if layer := packet.Layer(layers.LayerTypeIPv4); layer != nil {
				network := layer.(*layers.IPv4)
				ip = make([]byte, len(network.SrcIP))
				copy(ip, network.SrcIP)
			} else if layer := packet.Layer(layers.LayerTypeIPv6); layer != nil {
				network := layer.(*layers.IPv6)
				ip = make([]byte, len(network.SrcIP))
				copy(ip, network.SrcIP)
			}
			addr := &net.TCPAddr{IP: ip, Port: int(transport.SrcPort)}

			var init bool
			conn.lockflow(addr, func(e *tcpFlow) {
				e.ts = time.Now()
				e.seq = transport.Ack // update sequence number for every incoming packet
				if transport.SYN {    // for SYN packets, try initialize the flow entry once
					select {
					case <-e.ready:
					default:
						e.ack = transport.Seq + 1
						e.handle = handle

						// create link layer for WriteTo
						if layer := packet.Layer(layers.LayerTypeEthernet); layer != nil {
							ethLayer := layer.(*layers.Ethernet)
							e.linkLayer = &layers.Ethernet{
								EthernetType: ethLayer.EthernetType,
								SrcMAC:       ethLayer.DstMAC,
								DstMAC:       ethLayer.SrcMAC,
							}
						} else if layer := packet.Layer(layers.LayerTypeLoopback); layer != nil {
							loopLayer := layer.(*layers.Loopback)
							e.linkLayer = &layers.Loopback{Family: loopLayer.Family}
						}

						// create network layer for WriteTo
						if layer := packet.Layer(layers.LayerTypeIPv4); layer != nil {
							network := layer.(*layers.IPv4)
							e.networkLayer = &layers.IPv4{
								SrcIP:    network.DstIP,
								DstIP:    network.SrcIP,
								Protocol: network.Protocol,
								Version:  network.Version,
								Id:       network.Id,
								Flags:    layers.IPv4DontFragment,
								TTL:      64,
							}
						} else if layer := packet.Layer(layers.LayerTypeIPv6); layer != nil {
							network := layer.(*layers.IPv6)
							e.networkLayer = &layers.IPv6{
								Version:    network.Version,
								NextHeader: network.NextHeader,
								SrcIP:      network.DstIP,
								DstIP:      network.SrcIP,
								HopLimit:   64,
							}
						}

						// this tcp flow is ready to operate based on flow information
						if e.linkLayer != nil && e.networkLayer != nil {
							close(e.ready)
							init = true
						}
					}
				} else if transport.PSH {
					// Normal data push:
					// increase properly the ack number for other peer,
					// the other peer will update it's local sequence with the ack
					e.ack += uint32(len(transport.Payload))
					select {
					case conn.chMessage <- message{transport.Payload, addr}:
					case <-conn.die:
						return
					}
				} else if transport.FIN || transport.RST {
					e.ack++
					conn.closePeer(addr)
				}

			})
			if init { //send back SYN+ACk
				conn.writeToWithFlags(nil, addr, true, false, true, false, false)
			}
		}
	}()
}

// ReadFrom implements the PacketConn ReadFrom method.
func (conn *TCPConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	select {
	case <-conn.die:
		return 0, nil, io.EOF
	case packet := <-conn.chMessage:
		n = copy(p, packet.bts) // 用一次copy的方式，读取数据
		return n, packet.addr, nil
	}
}

// WriteTo implements the PacketConn WriteTo method.
func (conn *TCPConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return conn.writeToWithFlags(p, addr, false, true, true, false, false)
}

func (conn *TCPConn) writeToWithFlags(p []byte, addr net.Addr, SYN bool, PSH bool, ACK bool, FIN bool, RST bool) (n int, err error) {
	var ready chan struct{}
	conn.lockflow(addr, func(e *tcpFlow) { ready = e.ready })

	select {
	case <-conn.die:
		return 0, io.EOF
	case <-ready:
		tcpaddr, err := net.ResolveTCPAddr("tcp", addr.String())
		if err != nil {
			return 0, err
		}

		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}

		// fetch flow
		var flow tcpFlow
		conn.lockflow(addr, func(e *tcpFlow) { flow = *e })

		var localAddr *net.TCPAddr
		if conn.tcpconn != nil {
			localAddr = conn.tcpconn.LocalAddr().(*net.TCPAddr)
		} else {
			localAddr = conn.listener.Addr().(*net.TCPAddr)
		}
		tcp := &layers.TCP{
			SrcPort: layers.TCPPort(localAddr.Port),
			DstPort: layers.TCPPort(tcpaddr.Port),
			Window:  12580,
			Ack:     flow.ack,
			Seq:     flow.seq,
			SYN:     SYN,
			PSH:     PSH,
			ACK:     ACK,
			FIN:     FIN,
			RST:     RST,
		}

		tcp.SetNetworkLayerForChecksum(flow.networkLayer.(gopacket.NetworkLayer))

		payload := gopacket.Payload(p)

		gopacket.SerializeLayers(buf, opts, flow.linkLayer, flow.networkLayer, tcp, payload)
		if err := flow.handle.WritePacketData(buf.Bytes()); err != nil {
			return 0, err
		}

		conn.lockflow(addr, func(e *tcpFlow) { e.seq += uint32(len(p)) })
		return len(p), nil
	}
}

// Close closes the connection.
func (conn *TCPConn) Close() error {
	var err error
	conn.dieOnce.Do(func() {
		// signal connection has closed
		close(conn.die)

		// close all established tcp connections
		if conn.tcpconn != nil {
			conn.writeToWithFlags(nil, conn.tcpconn.RemoteAddr(), false, false, true, true, false)
			err = conn.tcpconn.Close()
		} else if conn.listener != nil {
			err = conn.listener.Close() // close listener
			conn.osConnsLock.Lock()
			for _, tcpconn := range conn.osConns { // close all accepted conns
				conn.writeToWithFlags(nil, tcpconn.RemoteAddr(), false, false, true, true, false)
				tcpconn.Close()
			}
			conn.osConns = nil
			conn.osConnsLock.Unlock()
		}

		// stop capturing
		for k := range conn.handles {
			conn.handles[k].Close()
		}
	})
	return err
}

// when a FIN or RST has arrived, trigger conn.Close on the original connection
// called from captureFlow
func (conn *TCPConn) closePeer(addr net.Addr) {
	if conn.tcpconn != nil {
		conn.tcpconn.Close()
	} else if conn.listener != nil {
		conn.osConnsLock.Lock()
		if c, ok := conn.osConns[addr.String()]; ok {
			c.Close()
			delete(conn.osConns, addr.String())
		}
		conn.osConnsLock.Unlock()
	}
}

// LocalAddr returns the local network address.
func (conn *TCPConn) LocalAddr() net.Addr {
	if conn.tcpconn != nil {
		return conn.tcpconn.LocalAddr()
	} else if conn.listener != nil {
		return conn.listener.Addr()
	}
	return nil
}

// SetDeadline implements the Conn SetDeadline method.
func (conn *TCPConn) SetDeadline(t time.Time) error { return errOpNotImplemented }

// SetReadDeadline implements the Conn SetReadDeadline method.
func (conn *TCPConn) SetReadDeadline(t time.Time) error { return errOpNotImplemented }

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (conn *TCPConn) SetWriteDeadline(t time.Time) error { return errOpNotImplemented }

// Dial connects to the remote TCP port,
// and returns a single packet-oriented connection
func Dial(network, address string) (*TCPConn, error) {
	// remote address resolve
	raddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	// create a dummy UDP socket, to get routing information
	dummy, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}

	// get iface name from the dummy connection, eg. eth0, lo0
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	var ifaceName string
	for _, iface := range ifaces {
		for _, addr := range iface.Addresses {
			if addr.IP.Equal(dummy.LocalAddr().(*net.UDPAddr).IP) {
				ifaceName = iface.Name
			}
		}
	}
	if ifaceName == "" {
		return nil, errors.New("cannot find correct interface")
	}

	// pcap init
	handle, err := pcap.OpenLive(ifaceName, 65536, false, time.Second)
	if err != nil {
		return nil, err
	}

	// TCP local address reuses the same address from UDP
	laddr, err := net.ResolveTCPAddr(network, dummy.LocalAddr().String())
	if err != nil {
		return nil, err
	}
	dummy.Close()

	// apply filter for incoming data
	filter := fmt.Sprintf("tcp and dst host %v and dst port %v and src host %v and src port %v", laddr.IP, laddr.Port, raddr.IP, raddr.Port)
	if err := handle.SetBPFFilter(filter); err != nil {
		return nil, err
	}

	// create an established tcp connection
	// will hack this tcp connection for packet transmission
	tcpconn, err := net.DialTCP(network, laddr, raddr)
	if err != nil {
		return nil, err
	}

	// fields
	conn := new(TCPConn)
	conn.die = make(chan struct{})
	conn.flowTable = make(map[string]tcpFlow)
	conn.handles = []*pcap.Handle{handle}
	conn.tcpconn = tcpconn
	conn.setTTL(tcpconn, 0) // prevent tcpconn from sending ACKs
	conn.chMessage = make(chan message)
	conn.captureFlow(handle)

	// discards data flow on tcp conn
	go func() {
		io.Copy(ioutil.Discard, tcpconn)
		tcpconn.Close()
	}()

	return conn, nil
}

// Listen acts like net.ListenTCP,
// and returns a single packet-oriented connection
func Listen(network, address string) (*TCPConn, error) {
	laddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	// get iface name from the dummy connection, eg. eth0, lo0
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	var handles []*pcap.Handle
	if laddr.IP == nil || laddr.IP.IsUnspecified() { // if address is not specified, capture on all ifaces
		for _, iface := range ifaces {
			if len(iface.Addresses) > 0 {
				// build dst host
				var dsthost = "("
				for k := range iface.Addresses {
					dsthost += "dst host " + iface.Addresses[k].IP.String()
					if k != len(iface.Addresses)-1 {
						dsthost += " or "
					} else {
						dsthost += ")"
					}
				}

				// try open on all nics
				if handle, err := pcap.OpenLive(iface.Name, 65536, false, time.Second); err == nil {
					// apply filter
					filter := fmt.Sprintf("tcp and %v and dst port %v", dsthost, laddr.Port)
					if err := handle.SetBPFFilter(filter); err != nil {
						return nil, err
					}

					handles = append(handles, handle)
				} else {
					return nil, err
				}
			}
		}
	} else {
		var ifaceName string
		for _, iface := range ifaces {
			for _, addr := range iface.Addresses {
				if addr.IP.Equal(laddr.IP) {
					ifaceName = iface.Name
				}
			}
		}
		if ifaceName == "" {
			return nil, errors.New("cannot find correct interface")
		}
		// pcap init
		handle, err := pcap.OpenLive(ifaceName, 65536, false, time.Second)
		if err != nil {
			return nil, err
		}

		// apply filter
		filter := fmt.Sprintf("tcp and dst host %v and dst port %v", laddr.IP, laddr.Port)
		if err := handle.SetBPFFilter(filter); err != nil {
			return nil, err
		}
		handles = []*pcap.Handle{handle}
	}

	// start listening
	l, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}

	// fields
	conn := new(TCPConn)
	conn.osConns = make(map[string]net.Conn)
	conn.handles = handles
	conn.flowTable = make(map[string]tcpFlow)
	conn.die = make(chan struct{})
	conn.listener = l
	conn.setTTL(l, 0)
	conn.chMessage = make(chan message)
	go conn.cleaner()

	for k := range handles {
		conn.captureFlow(handles[k])
	}

	// discard everything in original connection
	go func() {
		for {
			tcpconn, err := l.Accept()
			if err != nil {
				return
			}

			// record original connections for proper closing
			conn.osConnsLock.Lock()
			conn.osConns[tcpconn.LocalAddr().String()] = tcpconn
			conn.osConnsLock.Unlock()
			go func() { io.Copy(ioutil.Discard, tcpconn) }()
		}
	}()

	return conn, nil
}
