package src

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"golang.org/x/crypto/pbkdf2"
)

// TODO: create method for encoding and decoding values from and to base64

const ITERATIONS int = 10000
const KEYLEN int = 32
const SALTLEN int = 32
const JSON_PERM int = 0644
const VAULT_JSON string = "vault.json"
const VAULT_DIR string = "gpmanstore"

var ErrNotFound = errors.New("nothing found against provided site/service")
var ErrInvalidPaddingError = errors.New("crypto/padding: invalid padding size")
var ErrWrongPassphrase = errors.New("decryption Error, Passphrase might be wrong")

type UserPass struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func (us *UserPass) Encode(username, password, salt []byte) {
	us.Username = base64.StdEncoding.EncodeToString(username)
	us.Password = base64.StdEncoding.EncodeToString(password)
	us.Salt = base64.StdEncoding.EncodeToString(salt)
}

func (us UserPass) DecryptUserPass(passphrase string) (string, string, error) {
	var username string
	var password string
	salt_decoded, err := base64.StdEncoding.DecodeString(us.Salt)
	if err != nil {
		return username, password, err
	}

	username_decoded, err := base64.StdEncoding.DecodeString(us.Username)
	if err != nil {
		return username, password, err
	}

	password_decoded, err := base64.StdEncoding.DecodeString(us.Password)
	if err != nil {
		return username, password, err
	}

	username_dec, err := Decrypt(passphrase, salt_decoded, username_decoded)
	if err != nil {
		return username, password, err
	}
	password_dec, err := Decrypt(passphrase, salt_decoded, password_decoded)
	if err != nil {
		return username, password, err
	}
	username = string(username_dec)
	password = string(password_dec)

	return username, password, nil

}

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
	mode := ecb.NewECBDecrypter(block)
	text := make([]byte, len(ciphertext))
	mode.CryptBlocks(text, ciphertext)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	text, err = padder.Unpad(text)
	if errors.As(err, &ErrInvalidPaddingError) {
		return nil, ErrWrongPassphrase
	}
	if err != nil {
		return nil, err
	}
	return text, nil

}

func JsonWriter(passphrase, site, username, password string) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	VAULT_PATH := filepath.Join(home, VAULT_DIR)
	err = os.MkdirAll(VAULT_PATH, os.ModePerm)
	if err != nil {
		return err
	}

	username_enc, salt, err := Encrypt(passphrase, nil, []byte(username))
	if err != nil {
		return err
	}

	password_enc, _, err := Encrypt(passphrase, salt, []byte(password))
	if err != nil {
		return err
	}

	// username_base, password_base, salt_base := ToBase64(username_enc, password_enc, salt)

	userpass := UserPass{}
	userpass.Encode(username_enc, password_enc, salt)

	jsonObj := map[string]UserPass{
		site: userpass,
	}

	masrhaled, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}

	JSON_FILE := filepath.Join(VAULT_PATH, VAULT_JSON)

	_, err = os.Stat(JSON_FILE)
	if !os.IsNotExist(err) {
		content, err := ioutil.ReadFile(JSON_FILE)
		if err != nil {
			return err
		}
		unmarhaled := make(map[string]UserPass, 0)
		err = json.Unmarshal(content, &unmarhaled)
		if err != nil {
			return err
		}
		unmarhaled[site] = jsonObj[site]
		masrhaled, err = json.Marshal(unmarhaled)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(JSON_FILE, masrhaled, fs.FileMode(JSON_PERM))
	if err != nil {
		return err
	}
	return nil

}

func JsonReader(passphrase, site string) (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	JSON_FILE := filepath.Join(home, VAULT_DIR, VAULT_JSON)
	content, err := ioutil.ReadFile(JSON_FILE)
	if err != nil {
		return "", "", err
	}
	unmarsheled := make(map[string]UserPass, 0)
	err = json.Unmarshal(content, &unmarsheled)
	if err != nil {
		return "", "", err
	}
	data, ok := unmarsheled[site]
	if !ok {
		return "", "", ErrNotFound
	}

	username, password, err := data.DecryptUserPass(passphrase)
	if err != nil {
		return "", "", err
	}

	return username, password, nil

}

func JsonDelete(site string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	JSON_FILE := filepath.Join(home, VAULT_DIR, VAULT_JSON)
	content, err := ioutil.ReadFile(JSON_FILE)
	if err != nil {
		return err
	}
	unmarsheled := make(map[string]UserPass, 0)
	err = json.Unmarshal(content, &unmarsheled)
	if err != nil {
		return err
	}
	_, ok := unmarsheled[site]
	if !ok {
		return ErrNotFound
	}
	delete(unmarsheled, site)

	marshled, err := json.Marshal(unmarsheled)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(JSON_FILE, marshled, fs.FileMode(JSON_PERM))
	if err != nil {
		return err
	}
	return nil
}

func genRandNum(nums []int64) (int64, error) {
	bg := big.NewInt(int64(len(nums)))
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, err
	}
	return nums[n.Int64()], nil
}

func GenerateRandomPswd(length int, specialChars bool) (string, error) {
	var nums []int64
	var res string
	if !specialChars {
		for i := 65; i <= 90; i++ {
			nums = append(nums, int64(i))
		}
		for i := 97; i <= 122; i++ {
			nums = append(nums, int64(i))
		}
	} else {
		for i := 32; i <= 126; i++ {
			nums = append(nums, int64(i))
		}
	}

	for i := 0; i < length; i++ {
		num, err := genRandNum(nums)
		if err != nil {
			return res, err
		}
		res = res + string(rune(num))
	}
	return res, nil

}
