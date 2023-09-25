package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/neel-bp/gpman/src"
)

var subcommands = []string{
	"store",
	"get",
	"list",
	"delete",
	"help",
	"gitauth",
	"help",
	"pull",
}

func siteServiceOptions(prefix string) ([]string, error) {
	res := make([]string, 0)
	home, err := os.UserHomeDir()
	if err != nil {
		return res, err
	}
	JSON_FILE := filepath.Join(home, src.VAULT_DIR, src.VAULT_JSON)
	content, err := ioutil.ReadFile(JSON_FILE)
	if err != nil {
		return res, err
	}
	unmarsheled := make(map[string]interface{}, 0)
	err = json.Unmarshal(content, &unmarsheled)
	if err != nil {
		return res, err
	}

	for k := range unmarsheled {
		if strings.HasPrefix(strings.ToLower(k), strings.ToLower(prefix)) {
			res = append(res, k)
		}
	}
	return res, nil

}

func genericPredictor(list []string, prefix string) []string {
	if prefix == "" || strings.TrimSpace(prefix) == "" {
		return list
	}
	res := make([]string, 0)
	for _, v := range list {
		if strings.HasPrefix(strings.ToLower(v), strings.ToLower(prefix)) {
			res = append(res, v)
		}
	}
	return res
}

func printer(s []string) {
	for _, v := range s {
		fmt.Println(v)
	}
}

func main() {
	if len(os.Args) < 2 {
		return
	}

	// 0 is executable name
	// 1 is name of the program for which complete is trying to get completions of
	// 2 is prefix that is being written on term
	// 3 is the word that preceeds the prefix that is being written
	args := os.Args[1:]
	// no subcommand is given that means name of program and word preciding the prefix would be same
	if args[1] == args[3] {
		printer(genericPredictor(subcommands, args[2]))
		return
	}

}
