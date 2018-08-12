package zcrypt

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/box"
)

func EncryptToPubKey(pubkey, message string) (s string, err error) {
	var (
		b                []byte
		recipientKey     [32]byte
		nonce            [24]byte
		ephemeralPrivKey *[32]byte
		ephemeralPubKey  *[32]byte
	)
	// decode the public key
	if b, err = base64.StdEncoding.DecodeString(pubkey); err != nil { // 将pubkey进行base64
		return "", errors.Wrap(err, "couldn't decode public key")
	}

	if len(b) != 32 {
		return "", errors.New("secrets must be 32 bytes long")
	}

	copy(recipientKey[:], b[0:32]) // 只要32个字节

	if ephemeralPubKey, ephemeralPrivKey, err = box.GenerateKey(rand.Reader); err != nil {
		return "", errors.Wrap(err, "couldn't generate ephemeral key")
	}

	// read a random nonce
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", errors.Wrap(err, "couldn't read random nonce")
	}
	// encrypt the message, and append to the nonce
	b = box.Seal(nonce[:], []byte(message), &nonce, &recipientKey, ephemeralPrivKey)
	// prepend the ephemeral public key
	b = append((*ephemeralPubKey)[:], b...)
	// finally, return the encoded, encrypted message
	return base64.StdEncoding.EncodeToString(b), nil
}

// if m, err = zcrypt.EncryptToPubKey("MaRI8ibKsgg+QqvRPDPRrh8NbOR2nsB2Mk81ctU4KEE=", "this is a secret message"); err != nil {
//   return errors.Wrap(err, "couldn't encrypt the secret message")
// }
// log.Print(m)
