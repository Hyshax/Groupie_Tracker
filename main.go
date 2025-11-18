package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Reponse struct {
	Message string `json:"message"`
}

func GetArtists() {
	// 1. Envoyer la requête
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 2. Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 3. (Optionnel) Décoder le JSON
	var data Reponse
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

}

func GetDates() {
	// 1
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 2
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data Reponse
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
}

func GetLocations() {
	// 1
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 2
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var data Reponse
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
}

func GetRealations() {
	// 1
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 2
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var data Reponse
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
}
