package books

import (
	"context"
	"fmt"
	"log"
	"neolib/database"
	"neolib/types"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
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

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
		return err
	}
	imagePath := os.Getenv("IMAGE_PATH")

	var books []types.Books
	for rows.Next() {
		var book types.Books
		var bookExtension string
		if err := rows.Scan(&book.ISBN, &book.Title, &bookExtension); err != nil {
			log.Fatal("Unable to scan the row: ", err)
			return err
		}

		book.Path = fmt.Sprintf("%s%s/%s%s", imagePath, username, book.ISBN, bookExtension)
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
