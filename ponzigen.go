package main

import (
	"log"
	"os"
)

var files = []string{"ponzi"}

func main() {
	log.SetFlags(0)
	log.SetPrefix("ponzigen: ")

	path := os.Args[1]
	g, err := NewGenerator(path)
	if err != nil {
		panic(err)
	}
	g.SetPackage("models")
	err = g.Write(os.Stdout)
	if err != nil {
		panic(err)
	}
}
