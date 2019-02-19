package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
)

const (
	PKCS1 int = iota
	PKCS8
)

//GenRsaKey 生成RSA密钥
//默认使用PKCS1编码
//@return privatekey publickey error
func GenRsaKey(bits int, pkc ...int) ([]byte, []byte, error) {
	pkctype := PKCS1
	if pkc != nil && len(pkc) > 0 {
		pkctype = pkc[0]
	}
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	privatebuf := &bytes.Buffer{}
	var derStream []byte
	if pkctype == PKCS1 {
		derStream = x509.MarshalPKCS1PrivateKey(privateKey)
	} else {
		derStream = marshalPKCS8PrivateKey(privateKey)
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privatebuf, block)
	if err != nil {
		return nil, nil, err
	}
	publicbuf := &bytes.Buffer{}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	if err != nil {
		return nil, nil, err
	}
	err = pem.Encode(publicbuf, block)
	if err != nil {
		return nil, nil, err
	}
	return privatebuf.Bytes(), publicbuf.Bytes(), nil
}

//marshalPKCS8PrivateKey PKCS8编码
func marshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, err := asn1.Marshal(info)
	if err != nil {
		log.Panic(err.Error())
	}
	return k
}

//Encrypt RSA加密
func Encrypt(publickey []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	// 用公钥加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

//EncryptNoPadding 无填充模式加密
func EncryptNoPadding(publickey []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	crypted := bcpowmod(bin2int(string(data)), big.NewInt(int64(pub.E)), pub.N)
	rb := bcdechex(crypted)
	rb = padstr(rb)
	return hex.DecodeString(rb)
}

//Decrypt RSA解密
func Decrypt(privateKey []byte, data []byte, pkc ...int) ([]byte, error) {
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
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}

//DecryptNoPadding NO_PADDING模式解密
func DecryptNoPadding(privateKey *rsa.PrivateKey, data []byte) []byte {
	c := new(big.Int).SetBytes(data)
	return c.Exp(c, privateKey.D, privateKey.N).Bytes()
}

//DecodePrivateKey 解析私钥
func DecodePrivateKey(privateKey []byte, pkc ...int) (*rsa.PrivateKey, error) {
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
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}
	prkI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv = prkI.(*rsa.PrivateKey)
	return priv, err
}

//SignSHA512 RSA签名 SHA512withRSA
func SignSHA512(privateKey []byte, data []byte, pkc ...int) ([]byte, error) {
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
	hashed := sha512.Sum512(data)
	ret, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA512, hashed[:])
	if err != nil {
		return nil, errors.New("签名失败:" + err.Error())
	}
	return ret, nil
}

//SignSHA1 RSA签名 SHA1withRSA
func SignSHA1(privateKey []byte, data []byte, pkc ...int) ([]byte, error) {
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
	hashed := sha1.Sum(data)
	ret, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, hashed[:])
	if err != nil {
		return nil, errors.New("签名失败:" + err.Error())
	}
	return ret, nil
}

//VerifySHA512 RSA验证签名 SHA512withRSA
func VerifySHA512(publickey []byte, data []byte, vStr []byte) (bool, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return false, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	hashed := sha512.Sum512(data)
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA512, hashed[:], vStr)
	if err != nil {
		return false, errors.New("签名验证失败:" + err.Error())
	}
	return true, nil
}

//VerifySHA1 RSA验证签名 SHA1withRSA
func VerifySHA1(publickey []byte, data []byte, vStr []byte) (bool, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return false, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	hashed := sha1.Sum(data)
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed[:], vStr)
	if err != nil {
		return false, errors.New("签名验证失败:" + err.Error())
	}
	return true, nil
}
