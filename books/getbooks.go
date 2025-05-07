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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		log.Fatal("Username not found in context")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	category := r.URL.Query().Get("category")

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	rows, err := database.SelectBooks(ctx, categoryPtr, username)
	if err != nil {
		fmt.Println("Unable to get the books")
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Unable to get the columns")
		log.Fatal(err)
	}
	fmt.Println("Columns: ", columns)

	var books []types.Books
	for rows.Next() {
		var book types.Books
		if err := rows.Scan(&book.ISBN, &book.Title); err != nil {
			fmt.Println("Unable to scan the row")
			log.Fatal(err)
		}
		books = append(books, book)
	}

	tmpl, err := template.ParseFiles("templates/books.tmpl")
	if err != nil {
		fmt.Println("Unable to load the template")
		log.Fatal(err)
	}

	if err := tmpl.Execute(w, books); err != nil {
		fmt.Println("Unable to execute the template")
		log.Fatal(err)
	}

	// w.Write([]byte("Retrieved all books")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Retrieved all books in the system")
}
