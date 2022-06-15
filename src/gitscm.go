package src

import (
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const GITAUTHFILE string = ".gpmangit"

type GitInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
	Salt     string `json:"salt"`
}

func (g *GitInfo) Encode(username, password, url, salt []byte) {
	g.Username = base64.StdEncoding.EncodeToString(username)
	g.Password = base64.StdEncoding.EncodeToString(password)
	g.Url = base64.StdEncoding.EncodeToString(url)
	g.Salt = base64.StdEncoding.EncodeToString(salt)
}

func (g GitInfo) Decode() ([]byte, []byte, []byte, []byte, error) {
	username, err := base64.StdEncoding.DecodeString(g.Username)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	password, err := base64.StdEncoding.DecodeString(g.Password)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	salt, err := base64.StdEncoding.DecodeString(g.Salt)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	url, err := base64.StdEncoding.DecodeString(g.Url)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return username, password, url, salt, nil
}

func (g *GitInfo) Encrypt(passphrase, username, password, url string) error {
	username_enc, salt, err := Encrypt(passphrase, nil, []byte(username))
	if err != nil {
		return err
	}
	password_enc, _, err := Encrypt(passphrase, salt, []byte(password))
	if err != nil {
		return err
	}
	url_enc, _, err := Encrypt(passphrase, salt, []byte(url))
	if err != nil {
		return err
	}

	g.Encode(username_enc, password_enc, url_enc, salt)

	return nil
}

func (g GitInfo) Decrypt(passphrase string) (string, string, string, error) {
	var username, password, url string
	username_b, password_b, url_b, salt_b, err := g.Decode()
	if err != nil {
		return username, password, url, err
	}

	username_dec, err := Decrypt(passphrase, salt_b, username_b)
	if err != nil {
		return username, password, url, err
	}
	username = string(username_dec)

	password_dec, err := Decrypt(passphrase, salt_b, password_b)
	if err != nil {
		return username, password, url, err
	}
	password = string(password_dec)

	url_dec, err := Decrypt(passphrase, salt_b, url_b)
	if err != nil {
		return username, password, url, err
	}
	url = string(url_dec)

	return username, password, url, nil

}

func GitAuthInit(passphrase, url, username, access_token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	Obj := GitInfo{}
	err = Obj.Encrypt(passphrase, username, access_token, url)
	if err != nil {
		return err
	}
	marsheled, err := json.Marshal(Obj)
	if err != nil {
		return err
	}

	FILE := filepath.Join(home, GITAUTHFILE)
	err = ioutil.WriteFile(FILE, marsheled, fs.FileMode(JSON_PERM))
	if err != nil {
		return err
	}

	return nil

}
