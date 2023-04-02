package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/rekoda-project/rekoda/cmd"
)

func init() {
	go func() {
		http.ListenAndServe("localhost:8969", nil) // pprof
	}()
}

func main() {
	cmd.Execute()
}
