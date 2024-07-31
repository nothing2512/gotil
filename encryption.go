package gotil

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type Encryption struct {
	secret []byte
	iv     []byte
}

func DefaultEncryption() *Encryption {
	return NewEncryption("00000000000000000000000000000000", "1111111111111111")
}

func NewEncryption(secret, iv string) *Encryption {
	return &Encryption{[]byte(secret), []byte(iv)}
}

func (e *Encryption) Encrypt(data string) string {
	plainText := []byte(data)

	block, err := aes.NewCipher(e.secret)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	encryptStream := cipher.NewCTR(block, e.iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	ivHex := hex.EncodeToString(e.iv)
	encryptedDataHex := hex.EncodeToString(cipherText)

	return encryptedDataHex[len(ivHex):]
}

func (e *Encryption) Decrypt(data string) string {
	block, err := aes.NewCipher(e.secret)
	if err != nil {
		return ""
	}

	cipherText, err := hex.DecodeString(data)
	if err != nil {
		return ""
	}

	if block.BlockSize() != len(e.iv) {
		return ""
	}

	ctr := cipher.NewCTR(block, e.iv)
	plainText := make([]byte, len(cipherText))
	ctr.XORKeyStream(plainText, cipherText)

	return string(plainText)
}
