package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	PORT = ":8080"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	cleanPath = strings.Trim(cleanPath, "/")

	segments := strings.Split(cleanPath, "/")

	fmt.Printf("Request path segments: %v\n", segments)

	if len(segments) != 1 || segments[0] == "" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	action := segments[0]
	if(action != "echo") {
		http.Error(w, "Only 'echo' is allowed", http.StatusBadRequest)
		return
	}

	encode(w, http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getRoot)

	err := http.ListenAndServe(PORT, mux)
  if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
