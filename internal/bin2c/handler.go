package bin2c

import (
	"net/http"
	"strconv"

	"github.com/sebomancien/tools/internal/tmpl"
	"github.com/sebomancien/tools/pkg/converter"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, "body1.html", "Bin2C", nil)
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
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
