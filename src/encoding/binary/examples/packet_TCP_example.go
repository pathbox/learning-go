package packet

const (
	DEFAULE_HEADER           = "[**********]"
	DEFAULT_HEADER_LENGTH    = 12
	DEFAULT_SAVE_DATA_LENGTH = 4
)

type Packet struct {
	Header         string // 头部
	HeaderLengh    int32  //头部长度
	SaveDataLength int32  // 数据长度
	Data           []byte //数据
}

func (p *Packet) SetHeader(header string) *Packet {
	p.Header = header
	p.HeaderLengh = int32(len([]byte(header)))
	return p
}
func NewDefaultPacket(data []byte) *Packet {
	return &Packet{DEFAULE_HEADER, DEFAULT_HEADER_LENGTH, DEFAULT_SAVE_DATA_LENGTH, data}
}

func (p *Packet) Packet() []byte {
	return append(append([]byte(p.Header), p.IntToBytes(int32(len(p.Data)))...), p.Data...) // header + length + data
}

func (p *Packet) UnPacket(readerChannel chan []byte) []byte {
	dataLen := int32(len(p.Data))
	var i int32
	for i = 0; i < dataLen; i++ {
		//Termiate for loop when the remaining data is insufficient .
		if dataLen < i+p.HeaderLengh+p.SaveDataLength {
			break
		}

		if string(p.Data[i:i+p.HeaderLengh]) == p.Header {
			saveDataLenBeginIndex := i + p.HeaderLengh // 跳过头部数据长度后，就到这个包数据部分的开始索引
			actualDataLen := p.BytesToInt(p.Data[saveDataLenBeginIndex : saveDataLenBeginIndex+p.SaveDataLength])
			//The remaining data is less than one package
			if dataLen < i+p.HeaderLengh+p.SaveDataLength+actualDataLen {
				break
			}
			//Get a packet
			packageData := p.Data[saveDataLenBeginIndex+p.SaveDataLength : saveDataLenBeginIndex+p.SaveDataLength+actualDataLen]
			//send pacakge data to reader channel
			readerChannel <- packageData
			//get next package index
			i += p.HeaderLengh + p.SaveDataLength + actualDataLen - 1
		}
	}
	if i >= dataLen {
		return []byte{}
	}

	retrn p.Data[:] // 返回,实际拆包已经传到readerChannel <- packageData
}

func (p *Packet) IntToBytes(i int32) []byte {
	byteBuffer := bytes.NewBuffer([]byte{}) // new a buffer
	binary.Write(byteBuffer, binary.BigEndian, i) // int to bytes
	return byteBuffer.Bytes()
}

func (p *Packet) BytesToInt(data []byte) int32 {
	var val int32
	byteBuffer := bytes.NewBuffer(data)
	binary.Read(byteBuffer, binary.BigEndian, &val) // bytes to interface , here is int32 value
	return val
}

/*
Client:

dataPackage := packet.NewDefaultPacket([]byte(jsonString)).Packet()
Client.Write(dataPackage)

Server:
readerChannel := make(chan []byte, 1024)
//Store truncated data
remainBuffer := make([]byte, 0)
//read unpackage data from buffered channel
go func(reader chan []byte) {
		for {
				packageData := <-reader
				//....balabala....
		}
}(readerChannel)
  remainBuffer = packet.NewDefaultPacket(append(remainBuffer,recvData)).UnPacket(readerChannel)


*/