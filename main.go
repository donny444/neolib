package main

import (
	"fmt"
	"log"
	"neolib/auth"
	"neolib/books"
	"neolib/database"
	"net/http"
	"os"
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
	http.Handle("/", fs)

	signupHandler := http.HandlerFunc(auth.SignUp)
	http.Handle(fmt.Sprintf("%s/auth/signup", path), CorsMiddleware(signupHandler))

	signinHandler := http.HandlerFunc(auth.SignIn)
	http.Handle(fmt.Sprintf("%s/auth/signin", path), CorsMiddleware(signinHandler))

	booksHandler := http.HandlerFunc(handleBooks)
	http.Handle(fmt.Sprintf("%s/%s/", path, bookPath), CorsMiddleware(booksHandler))

	bookHandler := http.HandlerFunc(handleBook)
	http.Handle(fmt.Sprintf("%s/%s/{book}", path, bookPath), CorsMiddleware(bookHandler))

	libraryHandler := http.HandlerFunc(handleLibrary)
	http.Handle(fmt.Sprintf("%s/%s", path, libraryPath), CorsMiddleware(libraryHandler))
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books.GetBooks(w, r)
		return
	case http.MethodPost:
		books.CreateBook(w, r)
		return
	}
}

func handleBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books.GetBook(w, r)
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

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", wd)

	http.ListenAndServe(":5000", nil)
}
