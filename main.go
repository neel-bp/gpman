package main

import (
	"log"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.JsonReader("Graffhead2", "Baba")
	if err != nil {
		log.Fatal(err)
	}
}
