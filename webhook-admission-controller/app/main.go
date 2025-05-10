package main

import (
	Router "_gorestapi-k8s/router"
	"log"
	"net/http"
	_ "net/http/pprof" // automatically expose /debug/pprof/
	"os"
	// CPU Profiling: /debug/pprof/profile?seconds=30 (takes a 30-second snapshot)
	// Goroutines: /debug/pprof/goroutine
	// Heap (memory): /debug/pprof/heap
	// Block: /debug/pprof/block
	// OR
	// go tool pprof -http=localhost:6061 "http://localhost:6060/debug/pprof/profile?seconds=30"
	// Requires brew install graphviz
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		log.Println("Enabling pprof for profiling")
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	Router.Run()
}
