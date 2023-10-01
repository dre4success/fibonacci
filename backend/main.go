package main

import (
	"log"
	"net/http"

	"github.com/dre4success/fibonacci/fibonacci"
)

func serveFrontendApp() {
	const Port = ":8081"
	http.Handle("/", http.FileServer(http.Dir("./app/frontend/build")))
	log.Println("frontend app started on Port", Port)

	if err := http.ListenAndServe(Port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

func serveBackendApi() {
	const Port = ":8080"
	fibonacciServer := fibonacci.NewFibonacciServer()
	http.Handle("/api/", http.StripPrefix("/api", fibonacciServer))
	log.Println("backend api started on Port", Port)

	if err := http.ListenAndServe(Port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
func main() {
	go serveFrontendApp()
	serveBackendApi()
}
