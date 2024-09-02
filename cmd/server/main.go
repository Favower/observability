package main

import (
	"log"
	"net/http"

	"github.com/Favower/observability/internal/handlers"
	"github.com/Favower/observability/internal/storage"
)

func main() {
	storage := storage.NewMemStorage()

	http.HandleFunc("/update/", handlers.UpdateHandler(storage))

	log.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
