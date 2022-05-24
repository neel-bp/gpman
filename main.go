package main

import (
	"log"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.JsonReader("Graffhead2", "Ch")
	if err != nil {
		log.Fatal(err)
	}
}
