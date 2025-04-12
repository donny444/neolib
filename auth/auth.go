package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"neolib/database"

	"github.com/dgrijalva/jwt-go"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

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

	var credentials Credentials
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
