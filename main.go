package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexccamargo/wwuatt/services"
	_ "github.com/go-sql-driver/mysql"
)

// parsedTemplate is the global parsed HTML template.
// var parsedTemplate *template.Template

type movie struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func main() {
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// indexHandler handles requests to the / route.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		movies, err := getMoviesData()
		if err != nil {
			log.Printf("getMoviesData: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(movies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	// case "POST":
	// 	if err := saveVote(w, r); err != nil {
	// 		log.Printf("saveVote: %v", err)
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	}
	default:
		http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

// getMoviesData returns a templateData structure for populating the web page.
func getMoviesData() ([]movie, error) {
	var movies []movie

	rows, err := services.GetConnection().Query(`SELECT movie_id as ID, title FROM movies`)
	if err != nil {
		return []movie{}, fmt.Errorf("DB.QueryRow: %v", err)
	}

	for rows.Next() {
		nextMovie := movie{}
		err := rows.Scan(&nextMovie.ID, &nextMovie.Title)
		if err != nil {
			return []movie{}, fmt.Errorf("Rows.Scan: %v", err)
		}
		movies = append(movies, nextMovie)
	}

	return movies, nil
}
