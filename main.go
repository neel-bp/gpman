package main

import (
	"fmt"
	"os"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.CommandHandler(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}

}
