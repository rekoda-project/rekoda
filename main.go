package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/wmw64/rekoda/cmd"
)

func init() {
	go func() {
		http.ListenAndServe(":8969", nil) // pprof
	}()
}

func main() {
	cmd.Execute()
}
