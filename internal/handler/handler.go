package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sebomancien/bin2c/pkg/converter"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Bin2C</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
	</head>
	<body>
		<div class="container">
			<h1 class="mt-5">Binary to C Converter</h1>
			<form class="mt-3" action="/api/v1/bin2c" method="POST" enctype="multipart/form-data">
				<div class="mb-3">
					<label for="binary-file" class="form-label">Upload Binary File</label>
					<input class="form-control" type="file" id="binary-file" name="binary-file" required>
				</div>
				<div class="mb-3">
					<label for="array-name" class="form-label">C Array Name</label>
					<input class="form-control" type="text" id="array-name" name="array-name" placeholder="Enter C Array Name" required
					pattern="[a-zA-Z_][a-zA-Z0-9_]*" 
					title="C array name must start with a letter or underscore and can only contain letters, digits, and underscores">
					<div class="invalid-feedback">Invalid C array name. Must start with a letter or underscore and only contain alphanumeric characters or underscores.</div>
				</div>
				<div class="mb-3">
					<label for="bytes-per-line" class="form-label">Bytes per Line</label>
					<input class="form-control" type="number" id="bytes-per-line" name="bytes-per-line" placeholder="Enter Number of Bytes per Line" min="1" max="128" required>
					<div class="invalid-feedback">Bytes per line must be a number between 1 and 128.</div>
				</div>
				<button type="submit" class="btn btn-primary">Convert</button>
			</form>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, tmpl)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the request parameters
	arrayName := r.FormValue("array-name")
	if arrayName == "" {
		arrayName = converter.DefaultArrayName
	}
	bytesPerLineStr := r.FormValue("bytes-per-line")
	bytesPerLine, err := strconv.Atoi(bytesPerLineStr)
	if err != nil || bytesPerLine < 1 {
		bytesPerLine = converter.DefaultBytesPerLine
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
