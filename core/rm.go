package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type RMCharacters struct {
	Characters []RMCharacter
}

type RMCharacter struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func RetrieveCharacters(numberOfCharacters int) ([]RMCharacter, error) {
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
			return nil, e
		}
	}
	return characters, nil
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
