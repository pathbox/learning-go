package zcrypt

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/box"
)

func DecodeWithPrivKey(privKey, message string) (m string, err error) {
	var (
		b               []byte
		key             [32]byte
		nonce           [24]byte
		ephemeralPubKey [32]byte
		ok              bool
	)
	// decode the private key
	if b, err = base64.StdEncoding.DecodeString(privKey); err != nil {
		panic(err)
	}
	// if it's not exactly 32 bytes, it can't be a public key
	if len(b) != 32 {
		return "", errors.New("secrets must be 32 bytes long")
	}
	// turn the byte slice into a byte array
	copy(key[:], b[0:32])
	// decode the ciphertext
	if b, err = base64.StdEncoding.DecodeString(message); err != nil {
		return "", errors.Wrap(err, "couldn't decode encrypted message")
	}
	// extract the ephemeral public key
	if len(b) < 32 {
		return "", errors.New("missing ephemeral public key")
	}
	copy(ephemeralPubKey[:], b[:32])
	// extract the nonce
	if len(b) < 56 {
		return "", errors.New("missing nonce")
	}
	copy(nonce[:], b[32:56])
	// decrypt the message
	b, ok = box.Open(nil, b[56:], &nonce, &ephemeralPubKey, &key)
	if !ok {
		return "", errors.New("couldn't decrypt message")
	}
	return string(b), nil
}
