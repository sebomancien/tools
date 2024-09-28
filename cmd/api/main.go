package main

import (
	"net/http"
	"strconv"

	"github.com/sebomancien/bin2c/pkg/converter"
)

const (
	DEFAULT_ARRAY_NAME     = "myArray"
	DEFAULT_BYTES_PER_LINE = 32
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/bin2c", handler)
	http.ListenAndServe(":8080", router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get the request parameters
	arrayName := r.FormValue("array-name")
	if arrayName == "" {
		arrayName = DEFAULT_ARRAY_NAME
	}
	bytesPerLineStr := r.FormValue("bytes-per-line")
	bytesPerLine, err := strconv.Atoi(bytesPerLineStr)
	if err != nil || bytesPerLine < 1 {
		bytesPerLine = DEFAULT_BYTES_PER_LINE
	}

	// Read the binary data from the request
	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("binary-file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Write the c output
	w.Header().Set("Content-Disposition", "attachment; filename=\"output.c\"")
	w.Header().Set("Content-Type", "text/plain")

	config := converter.Config{
		ArrayName:   arrayName,
		BytePerLine: uint8(bytesPerLine),
	}

	converter.Convert(file, w, &config)
}
