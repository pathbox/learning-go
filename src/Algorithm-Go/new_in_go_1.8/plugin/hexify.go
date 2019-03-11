package main

import (
	"encoding/hex"
)

func Hexify(in string) string{
	return hex.EncodeToString([]byte(in))
}

//go build -buildmode=shared hexify.go


