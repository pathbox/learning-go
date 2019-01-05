package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"os"
)

func main() {
	var bits int
	flag.IntVar(&bits, "b", 2048, "密钥长度，默认为1024位")
	if err := GenRsaKey(bits); err != nil {
		log.Fatal("密钥文件生成失败！")
	}
	log.Println("密钥文件生成成功！")
}

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits) // rsa 加密算法
	if err != nil {
		return err
	}
	priStream := x509.MarshalPKCS1PrivateKey(privateKey) // 生成私钥byte stream
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: priStream,
	}
	priFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}

	defer priFile.Close()
	err = pem.Encode(priFile, block) // 将数据存到private.pem文件
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey                     // 公钥是在私钥的基础上生成的
	pubStream, err := x509.MarshalPKIXPublicKey(publicKey) // 生成公钥 byte stream
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLICK KEY",
		Bytes: pubStream,
	}
	pubFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, block)
	if err != nil {
		return err
	}
	return nil
}
