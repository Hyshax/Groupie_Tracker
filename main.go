package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// Artist correspond à un artiste renvoyé par l'API
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// GetArtists récupère la liste des artistes depuis l'API Groupie Tracker
func GetArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut HTTP inattendu : %s", resp.Status)
	}

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}

	return artists, nil
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

	// On sert les fichiers statiques depuis le dossier courant "."
	// de cette façon, /static/style.css correspondra à ./style.css
	fs := http.FileServer(http.Dir("."))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
