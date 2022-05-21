package src

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"golang.org/x/crypto/pbkdf2"
)

const ITERATIONS = 10000
const KEYLEN = 32
const SALTLEN = 32

func DeriveKey(passphrase string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, SALTLEN)
		_, err := rand.Read(salt)
		if err != nil {
			return nil, nil, err
		}
	}
	return pbkdf2.Key([]byte(passphrase), salt, ITERATIONS, KEYLEN, sha256.New), salt, nil
}

func Encrypt(passphrase string, salt, text []byte) ([]byte, []byte, error) {
	key, salt, err := DeriveKey(passphrase, salt)
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	text, err = padder.Pad([]byte(text))
	if err != nil {
		return nil, nil, err
	}
	ct := make([]byte, len(text))
	mode.CryptBlocks(ct, text)
	return ct, salt, nil
}

func Decrypt(passphrase string, salt, ciphertext []byte) ([]byte, error) {
	key, _, err := DeriveKey(passphrase, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := ecb.NewECBEncrypter(block)
	text := make([]byte, len(ciphertext))
	mode.CryptBlocks(text, ciphertext)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	text, err = padder.Unpad(text)
	if err != nil {
		return nil, err
	}
	return text, nil

}
