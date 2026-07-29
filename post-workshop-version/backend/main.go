package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Message struct {
	ID        int       `json:"id"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	rdb      *redis.Client
	ctx      = context.Background()
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for the workshop
		},
	}
	hub *Hub
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteJSON(message)
		}
	}
}

func main() {
	// Initialize Redis
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Verify connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	hub = newHub()
	go hub.run()

	http.HandleFunc("/messages", handleMessages)
	http.HandleFunc("/ws", serveWs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("Backend server starting on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan Message, 256)}
	client.hub.register <- client

	go client.writePump()
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

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		
		// Fetch messages from Redis
		msgStrings, err := rdb.LRange(ctx, "messages", 0, -1).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		messages := []Message{}
		for _, s := range msgStrings {
			var msg Message
			if err := json.Unmarshal([]byte(s), &msg); err != nil {
				continue
			}
			messages = append(messages, msg)
		}
		json.NewEncoder(w).Encode(messages)

	case http.MethodPost:
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Simple profanity filter
		forbiddenWords := []string{"fuck", "bitch", "shit", "asshole"}
		msg.Text = filterProfanity(msg.Text, forbiddenWords)

		// Get next ID from Redis
		id, err := rdb.Incr(ctx, "next_id").Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		msg.ID = int(id)
		msg.Timestamp = time.Now()

		// Store in Redis
		msgJSON, _ := json.Marshal(msg)
		if err := rdb.RPush(ctx, "messages", string(msgJSON)).Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Broadcast through the hub
		hub.broadcast <- msg

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func filterProfanity(text string, forbiddenWords []string) string {
	filtered := text
	for _, word := range forbiddenWords {
		replacement := strings.Repeat("*", len(word))
		// Use a simpler approach that ignores case but preserves the rest of the text
		// This is still a bit naive but better than ToLowering the whole thing
		lowerText := strings.ToLower(filtered)
		lowerWord := strings.ToLower(word)

		start := 0
		for {
			idx := strings.Index(lowerText[start:], lowerWord)
			if idx == -1 {
				break
			}
			idx += start
			filtered = filtered[:idx] + replacement + filtered[idx+len(word):]
			lowerText = lowerText[:idx] + replacement + lowerText[idx+len(word):]
			start = idx + len(word)
		}
	}
	return filtered
}
