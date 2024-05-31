package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

var database *sql.DB

const basePath = "/server"
const libraryPath = "library"
const bookPath = "book"

func SetupDatabase() {
	var err error
	database, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/neolib")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(database)
	database.SetConnMaxLifetime(time.Minute * 3)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)
}

func CorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization, Accept, Content-Length, Content-Type")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		handler.ServeHTTP(w, r)
	})
}

func SetupRoutes(path string) {
	bookHandler := http.HandlerFunc(handleBook)
	http.Handle(fmt.Sprintf("%s/%s/", path, bookPath), CorsMiddleware(bookHandler))
	libraryHandler := http.HandlerFunc(handleLibrary)
	http.Handle(fmt.Sprinf("%s/%s", path, libraryPath), CorsMiddleware(libraryHandler))
}

func main() {
	SetupDatabase()
	SetupRoutes(basePath)
	http.ListenAndServe(":5000", nil)
}
