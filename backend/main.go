package main

import (
	"net/http"

	"github.com/dre4success/fibonacci/fibonacci"
)

func main() {
	http.HandleFunc("/current", fibonacci.GetCurrent)
	http.HandleFunc("/next", fibonacci.GetNext)
	http.HandleFunc("/previous", fibonacci.GetPrevious)

	http.ListenAndServe(":8080", nil)
}
