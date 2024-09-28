package converter

import (
	"io"

	"github.com/sebomancien/bin2c/pkg/utils"
)

type Config struct {
	ArrayName   string
	BytePerLine uint8
}

func Convert(reader io.Reader, writer io.Writer, config *Config) {
	r := utils.NewReader(reader)
	w := utils.NewWriter(writer)

	w.Println("#include \"stdint.h\"")
	w.Println()
	w.Printf("uint8_t %s[] = {", config.ArrayName)
	w.Println()
	for line := range r.Chunk(int(config.BytePerLine)) {
		w.Print("    ")
		for _, data := range line {
			w.Printf("0x%02X, ", data)
		}
		w.Println()
	}
	w.Println("};")
}
