package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	messages = []Message{}
	mutex    sync.Mutex
	nextID   = 1
)

func main() {
	http.HandleFunc("/messages", handleMessages)

	port := ":8080"
	fmt.Printf("Backend server starting on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for frontend development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)

	case http.MethodPost:
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msg.ID = nextID
		nextID++
		msg.Timestamp = time.Now()
		messages = append(messages, msg)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
