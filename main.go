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

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

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
	artists := GetArtists()
	for _, artist := range artists {
		if artist.ID == id {
			return &artist
		}
	}
	return nil
}

func GetRelations(id int) *Relation {
	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching relations:", err)
		return nil
	}
	defer resp.Body.Close()

	var relation Relation
	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		log.Println("Error decoding relations:", err)
		return nil
	}
	return &relation
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	artists := GetArtists()
	tmpl.Execute(w, artists)
}

func apiArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists := GetArtists()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

func apiArtistByIDHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}
	artist := GetArtistByID(id)
	if artist == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

func apiRelationsHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	relation := GetRelations(id)
	if relation == nil {
		http.Error(w, "Relations introuvables", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relation)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/artists", enableCORS(apiArtistsHandler))
	http.HandleFunc("/api/artists/", enableCORS(apiArtistByIDHandler))

	http.HandleFunc("/api/relations/", enableCORS(apiRelationsHandler))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
