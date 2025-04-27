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

	err := validateSignUp(username, email, password)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.FindUser(ctx, username)
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
	w.Write([]byte("User created successfully"))
	fmt.Println("User created successfully")
}

func validateSignUp(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return fmt.Errorf("username, email, and password are required")
	}
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	if len(password) < 8 || len(password) > 16 {
		return fmt.Errorf("password must be between 8 and 16 characters")
	}
	if len(email) > 255 {
		return fmt.Errorf("email must be between 5 and 50 characters")
	}

	forbiddenChars := []rune{'!', '#', '$', '%', '^', '&', '*', '(', ')', '=', '+', '{', '}', '[', ']', '|', '\\', ':', ';', '\'', '"', '<', '>', ',', '?', '/', '`', '~', ' '}
	for _, char := range forbiddenChars {
		if containsRune(username, char) || containsRune(email, char) || containsRune(password, char) {
			return fmt.Errorf("username, email, and password must not contain forbidden special characters")
		}

	}

	return nil
}

func containsRune(s string, r rune) bool {
	for _, char := range s {
		if char == r {
			return true
		}
	}
	return false
}
