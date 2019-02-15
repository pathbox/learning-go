package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

//rsa长数据分段加解密

//EncodeData 加密,数据如果过长将使用分段加密
//	@param []byte => publickey 公钥数据
//	@param []byte => data 需要加密的内容
//	@param ...int => elen 数据分段加密长度,如果为空将按照密钥的能加密的最大长度分段
//	@return []byte 加密结果
//          error 错误信息
// 用公钥加密
func EncodeData(publickey []byte, data []byte, elen ...int) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	x := pub.N.BitLen()/8 - 11
	if elen != nil && len(elen) > 0 {
		if elen[0] > x {
			return nil, fmt.Errorf("指定长度数据过长,最大长度为 %d", x)
		}
		x = elen[0]
	}
	data = append(data, 1) //在数据最后加1,确保数据最后一个不是0,在解密的时候就可以直接把最后是0的数据去除因为那是padding的无效数据这样可以准确的得到有效的加密数据
	length := len(data)
	sb := bytes.NewBufferString("")
	for i := 0; i < length; {
		if length-i > x {
			ret, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data[i:i+x])
			i = i + x
			if err != nil {
				return nil, err
			}
			sb.Write(ret)
		} else {
			ret, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data[i:i+x])
			if err != nil {
				return nil, err
			}
			i = length
			sb.Write(ret)
		}
	}
	return sb.Bytes(), nil
}

//DecodeData 解密
//	@param []byte => privateKey 私钥
//	@param []byte => data 要解密的数据,如果数据进行过任何编码的请先解码
//	@param ...int => pkc 私钥的编码方式 pkcs1 pkcs8
//	@return []byte 解密内容
//          error 错误信息
// 用私钥解密
func DecodeData(privateKey []byte, data []byte, pkc ...int) ([]byte, error) {
	pkctype := PKCS1
	if pkc != nil && len(pkc) > 0 {
		pkctype = pkc[0]
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	var priv *rsa.PrivateKey
	var err error
	if pkctype == PKCS1 {
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else {
		prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv = prkI.(*rsa.PrivateKey)
	}
	if err != nil {
		return nil, err
	}
	length := len(data)
	x := priv.N.BitLen() / 8
	if length%x != 0 {
		return nil, errors.New("密文长度异常")
	}
	sb := &bytes.Buffer{}
	for i := 0; i < length; i = i + x {
		ret, err := rsa.DecryptPKCS1v15(rand.Reader, priv, data[i:i+x])
		if err != nil {
			return nil, err
		}
		sb.Write(ret)
	}
	ret := sb.Bytes()
	l := len(ret)
	if l < 1 {
		return ret, nil
	}
	i := 0
	for i < l {
		if ret[l-i-1] == 0 { //最后一个为0,继续判断
			i++
		} else {
			break
		}
	}
	return ret[:l-i-1], nil
}
