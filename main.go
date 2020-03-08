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

var service *services.MovieService

func main() {

	service = InitializeDI()
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
		movies, err := service.GetMoviesData()
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

	default:
		http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
	}
}
