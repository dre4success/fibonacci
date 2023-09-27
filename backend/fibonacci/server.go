package fibonacci

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var states sync.Map

type FibonacciState struct {
	Previous, Current int64
}

type JsonResponse struct {
	Token string `json:"token"`
	Value int64  `json:"value"`
}

func generateToken() (string, error) {
	buff := make([]byte, 16) // generate 128-bit token (16 bytes)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buff), nil
}

func respondWithJSON(w http.ResponseWriter, statusCode int, token string, value int64) {
	response := JsonResponse{Token: token, Value: value}
	jsonEncoded, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate a JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonEncoded)
}

func GetCurrent(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Fib-Token")
	if token == "" {
		var err error
		token, err = generateToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		log.Println("token: ", token)
		states.Store(token, &FibonacciState{Previous: 0, Current: 1})
	}
	stateInterface, ok := states.Load(token)
	if !ok {
		http.Error(w, "Invalid token provided", http.StatusBadRequest)
		return
	}
	// type checking the stateInterface that it is of FibonacciState type
	value := stateInterface.(*FibonacciState)
	respondWithJSON(w, http.StatusOK, token, value.Current)
}

func GetNext(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Fib-Token")
	if token == "" {
		http.Error(w, "Token not provided", http.StatusBadRequest)
		return
	}
	stateInterface, _ := states.Load(token)
	// type checking the stateInterface that it is of FibonacciState type
	value := stateInterface.(*FibonacciState)
	nextValue := value.Current + value.Previous
	value.Previous = value.Current
	value.Current = nextValue
	respondWithJSON(w, http.StatusOK, token, value.Current)
}

func GetPrevious(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Fib-Token")
	if token == "" {
		http.Error(w, "Token not provided", http.StatusBadRequest)
		return
	}
	stateInterface, _ := states.Load(token)
	// type checking the stateInterface that it is of FibonacciState type
	value := stateInterface.(*FibonacciState)
	prevValue := value.Current - value.Previous
	value.Current = value.Previous
	value.Previous = prevValue

	respondWithJSON(w, http.StatusOK, token, value.Current)
}
