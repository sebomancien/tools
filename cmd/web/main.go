package main

import (
	"fmt"
	"net/http"

	"github.com/sebomancien/bin2c/internal/handler"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/bin2c", handler.PostHandler)
	router.HandleFunc("GET /", handler.GetHandler)

	fmt.Println("Starting API server on port 8080")
	http.ListenAndServe(":8080", router)
}
