package main

import (
	"log"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.JsonDelete("Ch2")
	if err != nil {
		log.Fatal(err)
	}
}
