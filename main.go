package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	NWorkers = flag.Int("n", 256, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

const testMode bool = false

func main() {
	// Parse the command-line flags.
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	StartDispatcher(*NWorkers)

	mux := http.NewServeMux()

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	mux.HandleFunc("/work", Collector)

	if err := http.ListenAndServe(*HTTPAddr, mux); err != nil {
		fmt.Println(err.Error())
	}

}
