package rsa

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func bcdechex(decdata *big.Int) string {
	ret := ""
	m := big.NewInt(0)
	x := big.NewInt(16)
	for decdata.Int64() != 0 {
		m := m.Mod(decdata, x)
		decdata = decdata.Quo(decdata, x)
		ret = dechex(m) + ret
	}
	return ret
}

func bin2hex(bindata string) string {
	return hex.EncodeToString([]byte(bindata))
}

func bin2int(bindata string) *big.Int {
	bindata = bin2hex(bindata)
	return bchexdec(bindata)
}

func bcpowmod(num *big.Int, pw *big.Int, mod *big.Int) *big.Int {
	two := big.NewInt(2)
	ret := big.NewInt(1)
	tmp := big.NewInt(0)
	for {
		if 1 == tmp.Mod(pw, two).Int64() {
			ret = ret.Mul(ret, num)
			ret = ret.Mod(ret, mod)
		}
		num = num.Exp(num, two, mod)
		pw = pw.Quo(pw, two)
		if pw.Int64() == 0 {
			break
		}
	}
	return ret
}

func hexdec(hexStr string) float64 {
	data, _ := strconv.ParseInt(hexStr, 16, 64)
	return float64(data)
}

func bchexdec(hexdata string) *big.Int {
	len := len(hexdata)
	ret := big.NewInt(0)
	for i := range hexdata {
		dec := hexdec(hexdata[i : i+1])
		exp := len - i - 1
		pow := math.Pow(16.0, float64(exp))
		tmp, _ := big.NewInt(0).SetString(fmt.Sprintf("%.0f", pow*dec), 10)
		ret = ret.Add(ret, tmp)
	}
	return ret
}

func dechex(m *big.Int) string {
	return strconv.FormatInt(m.Int64(), 16)
}

func padstr(src string) string {
	src = strings.TrimSpace(src)
	padlen := 256 - len(src)
	if padlen > 0 {
		return strings.Repeat("0", padlen) + src
	}
	return src
}
