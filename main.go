package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
)

type RMCharacters struct {
	Characters []RMCharacter
}

type RMCharacter struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

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

func retrieveCharacter(resultChan chan<- *RMCharacter, errorChan chan<- error) {
	randomId := rand.Intn(826) + 1
	resp, err := http.Get(fmt.Sprintf("https://rickandmortyapi.com/api/character/%v", randomId))
	if err != nil {
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	rmCharacter := RMCharacter{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&rmCharacter)
	if err != nil {
		errorChan <- err
		return
	}

	resultChan <- &rmCharacter
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

		resultChan := make(chan *RMCharacter, numberOfCharacters)
		errorChan := make(chan error, numberOfCharacters)

		for range numberOfCharacters {
			go retrieveCharacter(resultChan, errorChan)
		}

		characters := make([]RMCharacter, numberOfCharacters)
		for i := range numberOfCharacters {
			select {
			case character := <-resultChan:
				characters[i] = *character
			case e := <-errorChan:
				slog.WarnContext(r.Context(), e.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		t.Execute(w, RMCharacters{characters})
	}
}
