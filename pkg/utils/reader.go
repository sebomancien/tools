package utils

import (
	"io"
	"iter"
	"log"
)

type Reader struct {
	reader io.Reader
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader: reader,
	}
}

func (r *Reader) Chunk(size int) iter.Seq[[]byte] {
	buffer := make([]byte, size)
	return func(yield func([]byte) bool) {
		for {
			n, err := r.reader.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				return
			}
			if !yield(buffer[:n]) {
				return
			}
		}
	}
}
