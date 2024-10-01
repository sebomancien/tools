package main

import (
	"fmt"
	"net/http"

	"github.com/sebomancien/tools/internal/bin2c"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /bin2c/{$}", bin2c.GetHandler)
	router.HandleFunc("POST /bin2c/convert", bin2c.ConvertHandler)

	fmt.Println("Starting API server on port http://localhost:8080/")
	http.ListenAndServe(":8080", router)
}
