// +build noserver,!noclient

package main

import (
	"io"
	"log"
)

const canServe = false

func serve(string, io.ReadWriteCloser) {
	log.Fatal("server mode is not available in this binary")
}
