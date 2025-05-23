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

func CorsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Accept", "multipart/form-data")
		w.Header().Add("Access-Control-Allow-Methods", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func SetupRoutes() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	signupHandler := http.HandlerFunc(auth.SignUp)
	http.Handle("/auth/signup/", CorsMiddleware(signupHandler))

	signinHandler := http.HandlerFunc(auth.SignIn)
	http.Handle("/auth/signin/", CorsMiddleware(signinHandler))

	booksHandler := http.HandlerFunc(handleBooks)
	http.Handle("/books/", CorsMiddleware(auth.Authentication(auth.Authorization(booksHandler))))

	bookHandler := http.HandlerFunc(handleBook)
	http.Handle("/books/{isbn}/", CorsMiddleware(auth.Authentication(auth.Authorization(bookHandler))))
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

func main() {
	err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database connected successfully")

	SetupRoutes()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", wd)

	http.ListenAndServe(":5000", nil)
}
