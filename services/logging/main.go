package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("LOGGING_PORT")
	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var logEntry map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Write to file if file logging is enabled
			if os.Getenv("LOG_OUTPUT") == "file" {
				logFile := "/var/log/micro/app.log"
				f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer f.Close()

				if err := json.NewEncoder(f).Encode(logEntry); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"status": "log stored"})
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Logging service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
