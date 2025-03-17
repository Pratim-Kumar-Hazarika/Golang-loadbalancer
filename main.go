package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	servers := []Server{
		newSimpleServer("http://localhost:9001"),
		newSimpleServer("http://localhost:9002"),
		newSimpleServer("http://localhost:9003"),
	}

	lb := NewLoadBalancer("8000", servers)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	})

	go startBackendServer(9001)
	go startBackendServer(9002)
	go startBackendServer(9003)

	fmt.Printf("Load balancer running at 'http://localhost:%s'\n", lb.port)
	err := http.ListenAndServe(":"+lb.port, nil)
	handleErr(err)
}

func startBackendServer(port int) {
	mux := http.NewServeMux() 

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 500) 
		fmt.Fprintf(w, "Response from backend server on port %d\n", port)
	})

	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting backend server on %s\n", serverAddr)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux, 
	}

	err := server.ListenAndServe()
	handleErr(err)
}
