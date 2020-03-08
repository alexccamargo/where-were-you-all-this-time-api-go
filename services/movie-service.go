package services

import (
	"fmt"

	"github.com/alexccamargo/wwuatt/models"
)

// MovieService describe MovieService instance
type MovieService struct {
	Database *Database
}

// NewMovieService generates a MovieService instance
func NewMovieService(database *Database) *MovieService {
	return &MovieService{Database: database}
}

// GetMoviesData returns all movies
func (ms *MovieService) GetMoviesData() ([]models.Movie, error) {
	var movies []models.Movie

	rows, err := ms.Database.Db.Query(`SELECT movie_id as ID, title FROM movies`)
	if err != nil {
		return []models.Movie{}, fmt.Errorf("DB.QueryRow: %v", err)
	}

	for rows.Next() {
		nextMovie := models.Movie{}
		err := rows.Scan(&nextMovie.ID, &nextMovie.Title)
		if err != nil {
			return []models.Movie{}, fmt.Errorf("Rows.Scan: %v", err)
		}
		movies = append(movies, nextMovie)
	}

	return movies, nil
}
