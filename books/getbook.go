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

func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetBook function get called")

	uuid := r.PathValue("book")
	fmt.Println("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := database.SelectBook(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}

	var book types.Books
	if err := row.Scan(&book.UUID, &book.Title, &book.Publisher, &book.Category, &book.Author, &book.Page, &book.Language, &book.PublicationYear, &book.ISBN); err != nil {
		http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
		fmt.Println("Scan error: ", err)
		return
	}

	tmpl, err := template.ParseFiles("templates/book.tmpl")
	if err != nil {
		http.Error(w, "Unable to load the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to load the template"))
		return
	}

	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, "Unable to execute the template", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to execute the template"))
		return
	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Retrieved the book specified by UUID"))
	result := fmt.Sprintf("Retrieved the book specified by the UUID: %s", uuid)
	fmt.Println(result)
}
