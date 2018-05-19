package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc     chan []byte
	wc     chan *Message
	reader Reader
	writer Writer
}

type ReadFromTail struct {
	inode uint64
	fd    *os.File
	path  string
}

type WriteToInfluxDB struct {
	batch      uint16
	retry      uint8
	influxConf *InfluxConf
}

type InfluxConf struct {
	Addr, Username, Password, Database, Precision string
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

type Monitor struct {
	listenPort string
	startTime  time.Time
	tpsSli     []int
	systemInfo SystemInfo
}

type SystemInfo struct {
	HandleLine   int       `json:"handleLine"`   // 已经处理的日志行数
	Tps          float64   `json:"tps"`          // 系统吞出量
	ReadChanLen  int       `json:"readChanLen"`  // input channel 长度
	WriteChanLen int       `json:"writeChanLen"` // output channel 长度
	RunTime      string    `json:"runTime"`      // 运行总时间
	ErrInfo      ErrorInfo `json:"errInfo"`      // 错误信息
}

type ErrorInfo struct {
	ReadErr    int `json:"readErr"`
	ProcessErr int `json:"processErr"`
	WriteErr   int `json:"writeErr"`
}

type TypeMonitor int

const (
	TypeHandleLine TypeMonitor = iota
	TypeReadErr
	TypeProcessErr
	TypeWriteErr
)

var (
	path, influxDsn, listenPort string
	processNum, writeNum        int
	TypeMonitorChan             = make(chan TypeMonitor, 200)
)

func NewReader(path string) (Reader, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &ReadFromTail{
		inode: stat.Ino,
		fd:    f,
		path:  path,
	}, nil
}

func NewWriter(influxDsn string) (Writer, error) {
	influxDsnSli := strings.Split(influxDsn, "@")
	if len(influxDsnSli) < 5 {
		return nil, errors.New("param influxDns err")
	}
	return &WriteToInfluxDB{
		batch: 50,
		retry: 3,
		influxConf: &InfluxConf{
			Addr:      influxDsnSli[0],
			Username:  influxDsnSli[1],
			Password:  influxDsnSli[2],
			Database:  influxDsnSli[3],
			Precision: influxDsnSli[4],
		},
	}, nil
}

func NewLogProcess(reader Reader, writer Writer) *LogProcess {
	return &LogProcess{
		rc:     make(chan []byte, 200),
		wc:     make(chan *Message, 200),
		reader: reader,
		writer: writer,
	}
}

func (l *LogProcess) Process() {
	/**
	'$remote_addr\t$http_x_forwarded_for\t$remote_user\t[$time_local]\t$scheme\t"$request"\t$status\t$body_bytes_sent\t"$http_referer"\t"$http_user_agent"\t"$gzip_ratio"\t$upstream_response_time\t$request_time'
	*/

	rep := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		TypeMonitorChan <- TypeHandleLine
		ret := rep.FindStringSubmatch(string(v))
		if len(ret) < 13 {
			log.Println("wrong input data:", v)
			TypeMonitorChan <- TypeProcessErr
			continue
		}

		timeLocal, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			TypeMonitorChan <- TypeProcessErr
			log.Println("time parse error:", err)
			continue
		}

		request := ret[6]
		requestSli := strings.Split(request, " ")
		if len(requestSli) < 3 {
			TypeMonitorChan <- TypeProcessErr
			log.Println("input request wrong:", request)
			continue
		}
		method := strings.TrimLeft(requestSli[0], "\"")
		u, err := url.Parse(requestSli[1])
		if err != nil {
			TypeMonitorChan <- TypeProcessErr
			log.Println("input url parse error:", err)
			continue
		}
		path := u.Path
		scheme := ret[5]
		status := ret[7]
		bytesSent, _ := strconv.Atoi(ret[8])
		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)

		l.wc <- &Message{
			TimeLocal:    timeLocal,
			Path:         path,
			Method:       method,
			Scheme:       scheme,
			Status:       status,
			BytesSent:    bytesSent,
			UpstreamTime: upstreamTime,
			RequestTime:  requestTime,
		}
	}
}

func (m *Monitor) start(lp *LogProcess) {
	// 消费监控数据
	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeHandleLine:
				m.systemInfo.HandleLine += 1
			case TypeReadErr:
				m.systemInfo.ErrInfo.ReadErr += 1
			case TypeProcessErr:
				m.systemInfo.ErrInfo.ProcessErr += 1
			case TypeWriteErr:
				m.systemInfo.ErrInfo.WriteErr += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	// 计算TPS
	go func() {
		for {
			<-ticker.C
			m.tpsSli = append(m.tpsSli, m.systemInfo.HandleLine)
			if len(m.tpsSli) > 2 {
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, m.systemStatus(lp))
	})
	log.Fatal(http.ListenAndServe(":"+m.listenPort, nil))
}

func (m *Monitor) systemStatus(lp *LogProcess) string {
	d := time.Now().Sub(m.startTime)
	m.systemInfo.RunTime = d.String()
	m.systemInfo.ReadChanLen = len(lp.rc)
	m.systemInfo.WriteChanLen = len(lp.wc)
	if len(m.tpsSli) >= 2 {
		// return math.Trunc(float64(m.tpsSli[1]-m.tpsSli[0])/5*1e3+0.5) * 1e-3
		m.systemInfo.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5
	}
	res, _ := json.MarshalIndent(m.systemInfo, "", "\t")
	return string(res)
}

func (r *ReadFromTail) Read(rc chan []byte) {
	defer close(rc)
	var stat syscall.Stat_t

	r.fd.Seek(0, 2) // seek 到末尾
	bf := bufio.NewReader(r.fd)

	for {
		line, err := bf.ReadBytes('\n')
		if err == io.EOF {
			if err := syscall.Stat(r.path, &stat); err != nil {
				// 文件切割，但新文件还没有生成
				time.Sleep(1 * time.Second)
			} else {
				nowInode := stat.Ino
				if nowInode == r.inode {
					// 无新的数据产生
					time.Sleep(1 * time.Second)
				} else {
					// 文件切割，重新打开文件
					r.fd.Close()
					fd, err := os.Open(r.path)
					if err != nil {
						panic(fmt.Sprintf("Open file err: %s", err.Error()))
					}
					r.fd = fd
					bf = bufio.NewReader(fd)
					r.inode = nowInode
				}
			}
			continue
		} else if err != nil {
			log.Printf("readFromTail ReadBytes err: %s", err.Error())
			TypeMonitorChan <- TypeReadErr
			continue
		}

		rc <- line[:len(line)-1]
	}
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	// https://github.com/influxdata/influxdb/tree/master/client
	infClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     w.influxConf.Addr,
		Username: w.influxConf.Username,
		Password: w.influxConf.Password,
	})
	if err != nil {
		panic(fmt.Sprintf("influxdb NewHTTPClient err:%s", err.Error()))
	}

	for {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  w.influxConf.Database,
			Precision: w.influxConf.Precision,
		})
		if err != nil {
			panic(fmt.Sprintf("influxdb NewBatchPoints error:%s", err.Error()))
		}

		var count uint16
	Fetch:
		for v := range wc {
			tags := map[string]string{
				"Path":   v.Path,
				"Scheme": v.Scheme,
				"Status": v.Status,
			}

			fields := map[string]interface{}{
				"UpstreamTime": v.UpstreamTime,
				"RequestTime":  v.RequestTime,
				"BytesSent":    v.BytesSent,
			}

			pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
			if err != nil {
				TypeMonitorChan <- TypeWriteErr
				log.Println("influxdb NewPoint error:", err)
				continue
			}

			bp.AddPoint(pt)
			count++
			if count > w.batch {
				break Fetch
			}
		}

		var i uint8
		for i = 1; i <= w.retry; i++ {
			if err := infClient.Write(bp); err != nil {
				TypeMonitorChan <- TypeWriteErr
				log.Printf("influxdb Write error:%s, retry:%d", err.Error(), i)
				time.Sleep(1 * time.Second)
			} else {
				log.Println(w.batch, "point has written")
				break
			}
		}
	}
}

func init() {
	flag.StringVar(&path, "path", "", "log file path")
	// influxDsn: http://ip:port@username@password@db@precision
	flag.StringVar(&influxDsn, "influxDsn", "", "influxDB dsn")
	flag.StringVar(&listenPort, "listenPort", "9193", "monitor port")
	flag.IntVar(&processNum, "processNum", 1, "process goroutine num")
	flag.IntVar(&writeNum, "writeNum", 1, "write goroutine num")
	flag.Parse()
}

func main() {
	reader, err := NewReader(path)
	if err != nil {
		panic(err)
	}

	writer, err := NewWriter(influxDsn)
	if err != nil {
		panic(err)
	}

	lp := NewLogProcess(reader, writer)

	go lp.reader.Read(lp.rc)

	for i := 1; i <= processNum; i++ {
		go lp.Process()
	}

	for i := 1; i <= writeNum; i++ {
		go lp.writer.Write(lp.wc)
	}

	m := &Monitor{
		listenPort: listenPort,
		startTime:  time.Now(),
	}
	go m.start(lp)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("capture exit signal:", s)
			os.Exit(1)
		case syscall.SIGUSR1: // 用户自定义信号
			log.Println(m.systemStatus(lp))
		default:
			log.Println("capture other signal:", s)
		}
	}
}
