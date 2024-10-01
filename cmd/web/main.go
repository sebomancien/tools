package main

import (
	"fmt"
	"net/http"

	"github.com/sebomancien/bin2c/internal/handler"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /bin2c/{$}", handler.GetHandler)
	router.HandleFunc("POST /bin2c/convert/{$}", handler.ConvertHandler)

	fmt.Println("Starting API server on port http://localhost:8080/")
	http.ListenAndServe(":8080", router)
}
