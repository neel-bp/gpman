package main

import (
	"log"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.JsonWriter("testman", "ch", "test1", "test2")
	if err != nil {
		log.Fatal(err)
	}
}
