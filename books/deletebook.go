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

	uuid := r.PathValue("book")
	fmt.Println("UUID: ", uuid)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DeleteBook(ctx, uuid)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to delete the book"))
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("The book specified by UUID is deleted"))
	result := fmt.Sprintf("Deleted the book specified by the UUID: %s", uuid)
	fmt.Println(result)
}
