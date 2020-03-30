package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// twitter snowflake epoch, 时间戳起始时间, 单位: 毫秒(ms)
	// 缺省值: 2019-03-22T18:30:00 +0800
	// 开始使用之后不要再修改，否则会导致ID重复
	// snowflake time = time.Now().UnixNano() / 1e6 - Epoch
	Epoch int64 = 1553248800000

	// 定义机器节点使用的位数
	// Machine + Step == 22
	MachineBits uint8 = 10

	// 定义自增ID使用的位数
	// Machine + Step == 22
	StepBits uint8 = 12

	machineMax   int64
	machineMask  int64
	stepMask     int64
	timeShift    uint8
	machineShift uint8
)

type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

// snowflake Node
type Node struct {
	mu      sync.Mutex
	time    int64
	machine int64
	step    int64
}

// snowflake ID
type ID int64

// NewNode 返回一个新的snowflake Node
func NewNode(machineID int64) (*Node, error) {
	node := new(Node)
	node.machine = machineID

	machineMax = 1<<MachineBits - 1
	machineMask = machineMax << StepBits
	stepMask = 1<<StepBits - 1
	timeShift = MachineBits + StepBits
	machineShift = StepBits

	if node.machine < 0 || node.machine > machineMask {
		return nil, errors.New("MachineID must be between 0 and " + strconv.FormatInt(machineMax, 10))
	}

	return node, nil
}

// 生成唯一ID
func (n *Node) Generate() ID {
	n.mu.Lock()

	// 纳秒时间戳转毫秒时间戳, 1e6 = int64(time.Millisecond)
	now := time.Now().UnixNano() / 1e6

	if n.time == now { // 当前时间与上次时间相同, step++
		n.step = (n.step + 1) & stepMask

		// step超出范围
		if n.step == 0 {
			// 等待1ms
			for now <= n.time {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else if n.time > now { // 如果机器时间回退, 例: 闰秒;时间同步
		// 等待时间达到上次的时间, 防止ID重复
		for now <= n.time {
			now = time.Now().UnixNano() / 1e6
		}
		n.step = 0
	} else { // 当前时间与上次时间不同, step归零
		n.step = 0
	}

	// 记录此次生成时间
	n.time = now

	// 通过位移把数据放到指定位置
	r := ID(
		(now-Epoch)<<timeShift |
			(n.machine << machineShift) |
			(n.step),
	)

	n.mu.Unlock()
	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}

// ID转换为10进制字符串
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// func (f ID) Base32() string {
// 	return intconv.Base32(uint64(f))
// }

// func ParseBase32(b []byte) (ID, error) {
// 	id, err := intconv.ParseBase32(string(b))
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ID(id), nil
// }

// func (f ID) Base58() string {
// 	return intconv.Base58(uint64(f))
// }

// func ParseBase58(b []byte) (ID, error) {
// 	id, err := intconv.ParseBase58(string(b))
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ID(id), nil
// }

// func (f ID) Base62() string {
// 	return intconv.Base62(uint64(f))
// }

// func ParseBase62(b []byte) (ID, error) {
// 	id, err := intconv.ParseBase62(string(b))
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ID(id), nil
// }

func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

func (f ID) Bytes() []byte {
	return []byte(f.String())
}

func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

func (f ID) Machine() int64 {
	return int64(f) & machineMask >> machineShift
}

func (f ID) Step() int64 {
	return int64(f) & stepMask
}

func (f ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(f), 10)
	buff = append(buff, '"')
	return buff, nil
}

func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*f = ID(i)
	return nil
}
