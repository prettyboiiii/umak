package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/prettyboiiii/umak/kamu"
	"github.com/robfig/cron/v3"
)

var (
	DIARY_NUMBER string
	CRON_TABS    []string
)

var (
	logger     = log.New(&logStore{}, "", log.LstdFlags)
	upgrader   = websocket.Upgrader{}
	logStoreMu sync.Mutex
	logs       []string
	clients    = make(map[*websocket.Conn]struct{})
)

type logStore struct{}

func (ls *logStore) Write(p []byte) (n int, err error) {
	logStoreMu.Lock()
	defer logStoreMu.Unlock()
	log.Print(string(p))
	logs = append(logs, string(p))
	for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, p)
	}
	return len(p), nil
}

func init() {
	DIARY_NUMBER = os.Getenv("DIARY_NUMBER")
	if DIARY_NUMBER == "" {
		log.Fatal("DIARY_NUMBER is missing")
	}
	crontabs := os.Getenv("CRON_TABS")
	CRON_TABS = strings.Split(crontabs, ",")
	if len(CRON_TABS) == 0 {
		log.Fatal("CRON_TABS are missing")
	}
}

func main() {
	logger.Println("Welcome to Umak!")

	kamuObj := kamu.New()
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(logger)))
	for _, crontab := range CRON_TABS {
		if _, err := c.AddFunc(crontab, func() {
			if err := kamuObj.GetPlaceInQueue(DIARY_NUMBER, 0); errors.Is(err, kamu.SeesionEndedErr) {
				if err := kamuObj.StartConversation(); err != nil {
					log.Fatal(err)
				}
				if err := kamuObj.GetPlaceInQueue(DIARY_NUMBER, 0); err != nil {
					log.Fatal(err)
				}
			} else if err != nil {
				log.Fatal(err)
			}
		}); err != nil {
			log.Fatal(err)
		}
	}

	c.Start()
	// Serve the log stream page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set response headers
		w.Header().Set("Content-Type", "text/html")

		// Serve the HTML page
		http.ServeFile(w, r, "./pages/index.html")
	})

	// Handle WebSocket connections
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		// Add the WebSocket connection to the clients map
		logStoreMu.Lock()
		clients[conn] = struct{}{}
		for _, logEntry := range logs {
			conn.WriteMessage(websocket.TextMessage, []byte(logEntry))
		}
		logStoreMu.Unlock()

		// Remove the WebSocket connection from the clients map when the connection is closed
		defer func() {
			logStoreMu.Lock()
			delete(clients, conn)
			logStoreMu.Unlock()
		}()

		// Keep the WebSocket connection open indefinitely
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	})
	// Start the HTTP server
	if err := http.ListenAndServe(":6900", nil); err != nil {
		log.Fatal(err)
	}

	// Keep the main function running
	// until the application is terminated
	select {}
}
