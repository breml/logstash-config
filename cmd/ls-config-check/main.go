package main

import (
	"log"
	"os"

	"github.com/breml/logstash-config"
)

func main() {
	in := os.Stdin
	nm := "stdin"
	log.SetFlags(0)
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		in = f
		nm = os.Args[1]
	}

	_, err := config.ParseReader(nm, in)
	if err != nil {
		errMsg := err.Error()
		if farthestFailure, hasErr := config.GetFarthestFailure(); hasErr {
			errMsg = farthestFailure
		}
		log.Fatal(errMsg)
	}
}
