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
}
