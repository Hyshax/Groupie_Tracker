package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := GetArtists()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		log.Println("GetArtists error:", err)
		return
	}

	if err := tmpl.Execute(w, artists); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		log.Println("template error:", err)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
