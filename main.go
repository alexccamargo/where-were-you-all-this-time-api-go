// https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// db is the global database connection pool.
var db *sql.DB

// parsedTemplate is the global parsed HTML template.
// var parsedTemplate *template.Template

// vote struct contains a single row from the votes table in the database.
// Each vote includes a candidate ("TABS" or "SPACES") and a timestamp.
type movie struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func main() {
	var err error

	// parsedTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("unable to parse template file: %s", err)
	}

	// If the optional DB_TCP_HOST environment variable is set, it contains
	// the IP address and port number of a TCP connection pool to be created,
	// such as "127.0.0.1:3306". If DB_TCP_HOST is not set, a Unix socket
	// connection pool will be created instead.
	if os.Getenv("DB_TCP_HOST") != "" {
		db, err = initTCPConnectionPool()
		if err != nil {
			log.Fatalf("initTcpConnectionPool: unable to connect: %s", err)
		}
	} else {
		db, err = initSocketConnectionPool()
		if err != nil {
			log.Fatalf("initSocketConnectionPool: unable to connect: %s", err)
		}
	}

	// Create the votes table if it does not already exist.
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS movies
	( movie_id SERIAL NOT NULL, title CHAR(255) NOT NULL, PRIMARY KEY (movie_id) );`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}
	// time_cast timestamp NOT NULL,

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

	rows, err := db.Query(`SELECT movie_id as ID, title FROM movies`)
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

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Printf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// initSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of MySQL.
func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPwd, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}

// initTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of MySQL.
func initTCPConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_tcp]
	var (
		dbUser    = mustGetenv("DB_USER")
		dbPwd     = mustGetenv("DB_PASS")
		dbTCPHost = mustGetenv("DB_TCP_HOST")
		dbName    = mustGetenv("DB_NAME")
	)

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPwd, dbTCPHost, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_tcp]
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]
}
