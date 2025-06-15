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

func GetBooks(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("GetBooks function get called")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		log.Fatal("Username not found in context")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	category := r.URL.Query().Get("category")
	title := r.URL.Query().Get("searchTerm")

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	var titlePtr *string
	if title != "" {
		titlePtr = &title
	}

	rows, err := database.SelectBooks(ctx, categoryPtr, titlePtr, username)
	if err != nil {
		log.Fatal("Unable to get the books: ", err)
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Unable to get the columns: ", err)
		return err
	}
	fmt.Println("Columns: ", columns)

	var books []types.Books
	for rows.Next() {
		var book types.Books
		if err := rows.Scan(&book.ISBN, &book.Title, &book.Path); err != nil {
			log.Fatal("Unable to scan the row: ", err)
			return err
		}
		books = append(books, book)
	}

	tmpl, err := template.ParseFiles("templates/books.tmpl")
	if err != nil {
		log.Fatal("Unable to load the template: ", err)
		return err
	}

	if err := tmpl.Execute(w, books); err != nil {
		log.Fatal("Unable to execute the template ", err)
		return err
	}

	// w.Write([]byte("Retrieved all books")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Retrieved all the books that passed the searching criteria")

	// return nil to indicate no error occurred
	return nil
}
