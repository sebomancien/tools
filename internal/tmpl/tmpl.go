package tmpl

import (
	"embed"
	"io"
	"log"
	"text/template"
)

//go:embed *.html
var files embed.FS
var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.New("").ParseFS(files, "*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func Execute(w io.Writer, content string, title string, data any) error {
	return tmpl.ExecuteTemplate(w, "base", struct {
		Title   string
		Content string
		Data    any
	}{
		Title:   title,
		Content: content,
		Data:    data,
	})
}
