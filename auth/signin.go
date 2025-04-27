package auth

import (
	"context"
	"fmt"
	"neolib/database"
	"neolib/types"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Println("Method not allowed")
		return
	}

	usernameOrEmail := r.FormValue("usernameOrEmail")
	password := r.FormValue("password")

	err := validateSignIn(usernameOrEmail, password)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := database.FindUser(ctx, usernameOrEmail)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Println("Invalid credentials")
		return
	}

	var credentials types.Credentials
	if err := user.Scan(&credentials.Username, &credentials.Email, &credentials.Password); err != nil {
		http.Error(w, "Unable to scan the row", http.StatusInternalServerError)
		fmt.Println("Unable to scan the row")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Println("Invalid credentials")
		return
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": credentials.Username,
		"email":    credentials.Email,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
	fmt.Println("User signed in successfully")
}

func validateSignIn(usernameOrEmail string, password string) error {
	if usernameOrEmail == "" {
		return fmt.Errorf("username or email is required")
	}

	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) > 16 {
		return fmt.Errorf("password must be at least 16 characters long")
	}

	forbiddenChars := []rune{'!', '#', '$', '%', '^', '&', '*', '(', ')', '=', '+', '{', '}', '[', ']', '|', '\\', ':', ';', '\'', '"', '<', '>', ',', '?', '/', '`', '~', ' '}
	for _, char := range forbiddenChars {
		if containsRune(usernameOrEmail, char) {
			return fmt.Errorf("username or email must not contain forbidden special characters")
		}

		if containsRune(password, char) {
			return fmt.Errorf("password must not contain forbidden special characters")
		}
	}

	return nil
}
