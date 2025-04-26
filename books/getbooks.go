package books

import (
	"context"
	"fmt"
	"log"
	"neolib/database"
	"neolib/types"
	"net/http"
	"text/template"
	"time"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetBooks function get called")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	category := r.URL.Query().Get("category")

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	rows, err := database.SelectBooks(ctx, categoryPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Unable to get the columns", http.StatusInternalServerError)
		return
	}
	fmt.Println("Columns: ", columns)

	var books []types.Books
	for rows.Next() {
		var book types.Books
		if err := rows.Scan(&book.ISBN, &book.Title); err != nil {
			http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
			fmt.Println("Scan error: ", err)
			return
		}
		books = append(books, book)
	}

	tmpl, err := template.ParseFiles("templates/books.tmpl")
	if err != nil {
		http.Error(w, "Unable to load the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to load the template"))
		return
	}

	if err := tmpl.Execute(w, books); err != nil {
		http.Error(w, "Unable to execute the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to execute the template"))
		return
	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Retrieved all books"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Println("Retrieved all books in the system")
}
