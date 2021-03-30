package main

import (
	"fmt"
	"math/big"
	"net"
)
/*
IP4与int64之间的转换
*/
//把Int64转换成IP4的的字符串形式
func Int64ToIp4String(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

//IP4地址转换为Int64，方便存储时，减少内存占用，提升性能
func ip4StringToInt64(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

//判断是否是IP地址，同时支持IP4,IP6
func IsIP(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}
	return true
}

func main() {
	v := int64(255255)
	fmt.Println(Int64ToIp4String(v))
}

// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF
type IP []byte

// IPv4 address a.b.c.d.
func IPv4(a, b, c, d byte) IP {
	p := make(IP, IPv6len)
	copy(p, v4InV6Prefix)
	p[12] = a
	p[13] = b
	p[14] = c
	p[15] = d
	return p
}

func dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i := 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10+int(s[i]-'0')
		if n >= big {
			return big, i, false
		}
	}
	if i == 0 {
		return 0,0,false
	}
	return n, i, true
}

func parseIPv4(s string) IP {
	var p [4]byte
	for i := 0; i < 4; i++ {
		if len(s) == 0 {
			return nil
		}

		if i > 0 {
			if s[0] != '.' {
				return nil
			}
			s = s[1:]
		}

		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return nil
		}
		s = s[c:]
		p[i] = byte(n)
	}

	if len(s) != 0 {
		return nil
	}
	return IPv4(p[0], p[1], p[2], p[3])
}

// To4 converts the IPv4 address ip to a 4-byte representation.
// If ip is not an IPv4 address, To4 returns nil.
func (ip IP) To4() IP {
	if len(ip) == IPv4len {
		return ip
	}
	if len(ip) == IPv6len &&
		isZeros(ip[0:10]) &&
		ip[10] == 0xff &&
		ip[11] == 0xff {
		return ip[12:16]
	}
	return nil
}