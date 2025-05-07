package books

import (
	"context"
	"fmt"
	"log"
	"neolib/database"
	"net/http"
	"time"
)

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteBook function get called")

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		log.Fatal("Username not found in context")
	}

	isbn := r.PathValue("isbn")
	fmt.Println("ISBN: ", isbn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DeleteBook(ctx, username, isbn)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to delete the book"))
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("The book specified by ISBN is deleted"))
	result := fmt.Sprintf("Deleted the book specified by the ISBN: %s", isbn)
	fmt.Println(result)
}
