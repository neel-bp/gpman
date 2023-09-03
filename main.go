package main

import (
	"fmt"
	"os"

	"github.com/neel-bp/gpman/src"
)

func main() {
	src.GenerateCompletion()
	err := src.CommandHandler(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
	}

}
