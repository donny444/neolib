package books

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"neolib/database"
	"net/http"
	"time"
)

func CheckBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CheckBook function get called")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		log.Fatal("Username not found in context")
	}

	isbn := r.PathValue("isbn")
	fmt.Println("ISBN: ", isbn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req struct {
		Check bool `json:"check"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal("Invalid JSON body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON body"))
		return
	}
	check := req.Check

	err := database.UpdateReadStatus(ctx, username, isbn, check)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to check the book"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book read status is updated"))
	result := fmt.Sprintf("Read status of the book with ISBN: %s is updated", isbn)
	fmt.Println(result)
}
