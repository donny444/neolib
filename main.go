package main

import (
	"fmt"
	"html/template"
	"log"
	"neolib/books"
	"neolib/database"
	"net/http"
	"os"
	"path/filepath"
)

const basePath = "/server"
const libraryPath = "library"
const bookPath = "books"

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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	bookHandler := http.HandlerFunc(handleBook)
	http.Handle(fmt.Sprintf("%s/%s/", path, bookPath), CorsMiddleware(bookHandler))
	libraryHandler := http.HandlerFunc(handleLibrary)
	http.Handle(fmt.Sprintf("%s/%s", path, libraryPath), CorsMiddleware(libraryHandler))
	http.Handle("/test", CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Test struct {
			Name   string
			Number int
		}
		var test Test
		test.Name = "Hello"
		test.Number = 1

		tmpl, err := template.ParseFiles("templates/test.tmpl")
		if err != nil {
			http.Error(w, "Unable to load the template", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, test); err != nil {
			http.Error(w, "Unable to execute the template", http.StatusInternalServerError)
			return
		}
	})))
	http.HandleFunc("/", serveTemplate)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Print(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func handleBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books.GetBooks(w, r)
		return
	case http.MethodPost:
		books.CreateBook(w, r)
		return
	case http.MethodPut:
		books.EditBook(w, r)
		return
	case http.MethodDelete:
		books.DeleteBook(w, r)
		return
	}
}

func handleLibrary(w http.ResponseWriter, r *http.Request) {
}

// func getBooksbyCategory(_ http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	books, err := database.QueryContext(ctx, "SELECT * FROM books WHERE category = ?", r.URL.Query().Get("category"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer books.Close()
// 	fmt.Println(books)
// }

func main() {
	database.SetupDatabase()
	SetupRoutes(basePath)
	http.ListenAndServe(":5000", nil)
}
