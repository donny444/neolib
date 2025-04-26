package books

import (
	"context"
	"fmt"
	"log"
	"neolib/database"
	"net/http"
	"time"
)

func EditBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("EditBook function get called")

	isbn := r.PathValue("isbn")
	fmt.Println("ISBN: ", isbn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Helper function to get pointer to string or nil
	optionalInput := func(value string) *string {
		if value == "" {
			return nil
		}
		return &value
	}

	requiredInput := func(value string) string {
		if value == "" {
			http.Error(w, "Required inputs are not filled", http.StatusBadRequest)
		}
		return value
	}

	err := database.UpdateBook(ctx,
		requiredInput(r.FormValue("isbn")),
		requiredInput(r.FormValue("title")),
		optionalInput(r.FormValue("publisher")),
		optionalInput(r.FormValue("category")),
		optionalInput(r.FormValue("author")),
		optionalInput(r.FormValue("page")),
		optionalInput(r.FormValue("language")),
		optionalInput(r.FormValue("publication_year")))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to update the book"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated"))
	result := fmt.Sprintf("Book with ISBN: %s is updated", isbn)
	fmt.Println(result)
}
