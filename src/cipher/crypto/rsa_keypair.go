package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func main() {
  priv, err := rsa.GenerateKey(rand.Reader, 2016)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(priv)
  err = priv.Validate()
  if err != nil {
    fmt.Println("Validation failed.", err)
  }

  priv_der := x509.MarshalPKCS1PrivateKey(priv)

  priv_blk := pem.Block{
    Type: "RSA PRIVATE KEY",
    Headers: nil,
    Bytes: priv_der,
  }

  priv_pem := string(pem.EncodeToMemory(&priv_blk))
  fmt.Printf(priv_pem)

  pub := priv.PublicKey
  pub_der, err := x509.MarshalPKIXPublicKey(&pub)
  if err != nil {
    fmt.Println("Failed to get der format for PublicKey.", err)
		return
  }

  pub_blk := pem.Block{
    Type: "PUBLIC KEY",
    Headers: nil,
    Bytes: pub_der,
  }
  pub_pem := string(pem.EncodeToMemory(&pub_blk))
  fmt.Printf(pub_pem)
}
