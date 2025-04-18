package auth

import (
	"context"
	"fmt"
	"neolib/database"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.FindUser(ctx, username)
	if err != nil {
		http.Error(w, "Username already used", http.StatusConflict)
		fmt.Println("Username already used")
		return
	}

	_, err = database.FindUser(ctx, email)
	if err != nil {
		http.Error(w, "Email already exists", http.StatusInternalServerError)
		fmt.Println("Email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		fmt.Println("Unable to hash password")
		return
	}

	err = database.CreateUser(ctx, username, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		fmt.Println("Unable to create user")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
