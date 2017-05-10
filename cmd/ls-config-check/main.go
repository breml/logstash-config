package main

import (
	"fmt"
	"log"
	"os"

	"github.com/breml/logstash-config"
)

func main() {
	in := os.Stdin
	nm := "stdin"
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
		nm = os.Args[1]
	}

	got, err := config.ParseReader(nm, in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(got)
}
