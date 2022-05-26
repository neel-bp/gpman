package main

import (
	"log"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.JsonReader("testman", "Ch1")
	if err != nil {
		log.Fatal(err)
	}
}
