package main

import (
	"log"

	"github.com/neel-bp/gpman/src"
)

func main() {
	err := src.JsonWriter("Graffhead2", "4chan", "nigger", "nigger")
	if err != nil {
		log.Fatal(err)
	}

}
