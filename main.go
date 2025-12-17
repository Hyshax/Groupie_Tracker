package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// Middleware CORS pour permettre les requêtes cross-origin
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func GetArtists() []Artist {
	// la fonction recupère les artistes depuis l'api et les decodes en struct Artist
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println("Error fetching artists:", err)
		return []Artist{}
	}
	defer resp.Body.Close()

	var artists []Artist
	json.NewDecoder(resp.Body).Decode(&artists)

	return artists
}

func GetArtistByID(id int) *Artist {
	// permet de recupérer un artiste par son ID
	artists := GetArtists()

	for _, artist := range artists {
		if artist.ID == id {
			return &artist
		}
	}

	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	artists := GetArtists()
	tmpl.Execute(w, artists)
}

// API endpoint pour récupérer tous les artistes en JSON
func apiArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists := GetArtists()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artists)
}

// API endpoint pour récupérer un artiste spécifique par ID
func apiArtistByIDHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	id, _ := strconv.Atoi(pathParts[2])
	artist := GetArtistByID(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artist)
}

func main() {
	// Routes HTML
	http.HandleFunc("/", homeHandler)

	// Routes API avec CORS
	http.HandleFunc("/api/artists", enableCORS(apiArtistsHandler))
	http.HandleFunc("/api/artists/", enableCORS(apiArtistByIDHandler))

	// Fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur sur http://localhost:8080")
	log.Println("API disponible sur:")
	log.Println("  - GET /api/artists (tous les artistes)")
	log.Println("  - GET /api/artists/{id} (artiste spécifique)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
