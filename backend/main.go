package main

import (
	"log"
	"net/http"

	"github.com/dre4success/fibonacci/fibonacci"
)


func main() {
	const Port = ":8080"

	fibonacciServer := fibonacci.NewFibonacciServer()

	http.Handle("/api/", http.StripPrefix("/api", fibonacciServer))

	log.Println("server started on Port", Port)
	if err := http.ListenAndServe(Port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
