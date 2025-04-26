package books

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"neolib/database"
	"neolib/types"
	"net/http"
	"time"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetBook function get called")

	w.Header().Add("Content-Type", "application/json")

	isbn := r.PathValue("isbn")
	fmt.Println("ISBN: ", isbn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := database.SelectBook(ctx, isbn)
	if err != nil {
		log.Fatal(err)
	}

	var book types.Books
	if err := row.Scan(&book.ISBN, &book.Title, &book.Publisher, &book.Category, &book.Author, &book.Pages, &book.Language, &book.PublicationYear); err != nil {
		http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
		fmt.Println("Scan error: ", err)
		return
	}

	response := map[string]interface{}{
		"isbn":            book.ISBN,
		"title":           book.Title,
		"publisher":       book.Publisher,
		"category":        book.Category,
		"author":          book.Author,
		"pages":           book.Pages,
		"language":        book.Language,
		"publicationYear": book.PublicationYear,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Unable to marshal response to JSON", http.StatusInternalServerError)
		log.Fatal("JSON marshal error: ", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	result := fmt.Sprintf("Retrieved the book specified by the ISBN: %s", isbn)
	fmt.Println(result)
}
