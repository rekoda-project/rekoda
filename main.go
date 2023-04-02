package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/kappaflow/rekoda/cmd"
)

func init() {
	go func() {
		http.ListenAndServe("localhost:8969", nil) // pprof
	}()
}

func main() {
	cmd.Execute()
}
