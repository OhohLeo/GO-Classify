package main

import (
	"log"
	"testing"
)

func TestReadDirectory(t *testing.T) {

	c := ReadDirectory(`/home/ohohleo`, false)

	for {

		if f, ok := <-c; ok {
			log.Printf("Received %s\n", f.FullPath)
			continue
		}

		break
	}
}
