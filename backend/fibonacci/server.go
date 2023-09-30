package fibonacci

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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

type FibonacciServer struct {
	http.Handler
}

func allowCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Fib-Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func NewFibonacciServer() *FibonacciServer {
	s := new(FibonacciServer)

	router := http.NewServeMux()
	router.HandleFunc("/current", allowCORS(s.getCurrent))
	router.HandleFunc("/next", allowCORS(s.getNext))
	router.HandleFunc("/previous", allowCORS(s.getPrevious))

	s.Handler = router
	return s
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

func (s *FibonacciServer) getCurrent(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Fib-Token")
	if token == "" {
		var err error
		token, err = generateToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		states.Store(token, &FibonacciState{Current: 0, Previous: 1})
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

func (s *FibonacciServer) getNext(w http.ResponseWriter, r *http.Request) {
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

func (s *FibonacciServer) getPrevious(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Fib-Token")
	if token == "" {
		http.Error(w, "Token not provided", http.StatusBadRequest)
		return
	}
	stateInterface, _ := states.Load(token)
	// type checking the stateInterface that it is of FibonacciState type
	value := stateInterface.(*FibonacciState)
	if value.Current == 0 {

		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"message": "Can't go back any further in the sequence.",
		}
		jsonResponse, _ := json.Marshal(response)
		w.Write(jsonResponse)
		return
	}
	newPrevious := value.Current - value.Previous
	value.Current = value.Previous
	value.Previous = newPrevious

	respondWithJSON(w, http.StatusOK, token, value.Current)
}
