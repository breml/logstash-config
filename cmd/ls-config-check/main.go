package main

import (
	"log"
	"os"

	"github.com/breml/logstash-config"
)

func checkFile(f *os.File, name string) bool {
	_, err := config.ParseReader(name, f)
	if err != nil {
		errMsg := err.Error()
		if farthestFailure, hasErr := config.GetFarthestFailure(); hasErr {
			errMsg = farthestFailure
		}
		log.Printf("%s: %s", name, errMsg)
		return false
	}
	return true
}

func main() {
	success := true
	log.SetFlags(0)
	if len(os.Args) > 1 {
		for _, path := range os.Args[1:] {
			f, err := os.Open(path)
			if err != nil {
				log.Print(err)
				success = false
				continue
			}
			success = checkFile(f, path) && success
			if err := f.Close(); err != nil {
				log.Print(err)
				success = false
			}
		}
	} else {
		success = checkFile(os.Stdin, "stdin") && success
	}
	if !success {
		os.Exit(1)
	}
}
