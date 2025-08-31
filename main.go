package main

import (
	"fmt"
	"log"
	"neolib/advanced"
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
	bookImages := http.StripPrefix("/images/", http.FileServer(http.Dir("images")))
	http.Handle("/images/", bookImages)

	signupHandler := http.HandlerFunc(auth.SignUp)
	http.Handle("/auth/signup/", CorsMiddleware(signupHandler))

	signinHandler := http.HandlerFunc(auth.SignIn)
	http.Handle("/auth/signin/", CorsMiddleware(signinHandler))

	booksHandler := http.HandlerFunc(handleBooks)
	http.Handle("/books/", CorsMiddleware(auth.Authentication(auth.Authorization(booksHandler))))

	bookHandler := http.HandlerFunc(handleBook)
	http.Handle("/books/{isbn}/", CorsMiddleware(auth.Authentication(auth.Authorization(bookHandler))))

	checkHandler := http.HandlerFunc(books.CheckBook)
	http.Handle("/books/{isbn}/check/", CorsMiddleware(auth.Authentication(auth.Authorization(checkHandler))))

	advancedHandler := http.HandlerFunc(handleAdvanced)
	http.Handle("/advanced/{insight}", CorsMiddleware(auth.Authentication(auth.Authorization(advancedHandler))))
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
		// if r.PathValue("option") == "check" {
		// 	books.CheckBook(w, r)
		// 	return
		// }
		books.EditBook(w, r)
		return
	case http.MethodDelete:
		books.DeleteBook(w, r)
		return
	}
}

func handleAdvanced(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.PathValue("insight") == "" {
		http.Error(w, "Missing insight parameter", http.StatusBadRequest)
		return
	}

	switch r.PathValue("insight") {
	case "top-five-categories":
		advanced.TopFiveCategories(w, r)
		return
	case "reading-status":
		advanced.ReadingStatusByCategory(w, r)
		return
		// case "books-by-month":
		// 	advanced.ReadBooksByMonth(w, r)
		// 	return
		// case "language-portions":
		// 	advanced.GetLanguagePortions(w, r)
		// 	return
		// case "read-counts":
		// 	advanced.GetReadCounts(w, r)
		// 	return
		// }
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
