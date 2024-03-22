package main

import (
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gossie/meetup/core"
)

//go:embed templates/*
var htmlTemplates embed.FS

func main() {
	customizeLogging()

	t, err := template.ParseFS(htmlTemplates, "templates/*.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("GET /random-characters", trace(profile(getCharacters(t))))

	slog.Info("started server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCharacters(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "retrieving random character")

		numberOfCharacters, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			numberOfCharacters = 1
		}

		if numberOfCharacters > 50 {
			numberOfCharacters = 50
		}

		characters, err := core.RetrieveCharacters(numberOfCharacters)
		if err != nil {
			slog.WarnContext(r.Context(), err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t.Execute(w, core.RMCharacters{Characters: characters})
	}
}
