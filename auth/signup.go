package auth

import (
	"context"
	"fmt"
	"log"
	"neolib/database"
	"net/http"
	"os"
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

	_, err = database.SelectUser(ctx, username, email)
	if err != nil {
		http.Error(w, "Username or email already been used", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Unable to hash password: ", err)
	}

	err = database.InsertUser(ctx, username, email, string(hashedPassword))
	if err != nil {
		log.Fatal("Unable to create user: ", err)
	}

	// Create a directory for the user at /images/[username] path
	userDir := fmt.Sprintf("./images/%s", username)
	err = os.Mkdir(userDir, 0755)
	if err != nil {
		log.Fatal("Unable to create user directory: ", err)
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
