package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupDatabase() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/neolib")
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func InsertBook(ctx context.Context, uuid string, title string, isbn string, publisher *string, category *string, author *string, page *string, language *string, publicationYear *string, fileContent []byte) error {
	_, err := db.ExecContext(ctx, "INSERT INTO books (uuid, title, publisher, category, author, page, language, publication_year, isbn, file_content) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		uuid,
		title,
		publisher,
		category,
		author,
		page,
		language,
		publicationYear,
		isbn,
		fileContent)
	return err
}

func SelectBooks(ctx context.Context) (*sql.Rows, error) {
	return db.QueryContext(ctx, "SELECT uuid, title, isbn FROM books")
}

func SelectBook(ctx context.Context, uuid string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT uuid, title, publisher, category, author, page, language, publication_year, isbn FROM books WHERE uuid = ?", uuid), nil
}

func UpdateBook(ctx context.Context, uuid string, title string, isbn string, publisher *string, category *string, author *string, page *string, language *string, publicationYear *string) error {
	_, err := db.ExecContext(ctx, "UPDATE books SET title = ?, isbn = ?, publisher = ?, category = ?, author = ?, page = ?, language = ?, publication_year = ? WHERE uuid = ?",
		title,
		isbn,
		publisher,
		category,
		author,
		page,
		language,
		publicationYear,
		uuid)
	return err
}

func DeleteBook(ctx context.Context, uuid string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM books WHERE uuid = ?", uuid)
	return err
}

func FindUser(ctx context.Context, usernameOrEmail string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT username, email, password FROM users WHERE username = ? OR email = ?", usernameOrEmail, usernameOrEmail), nil
}

func CreateUser(ctx context.Context, username string, email string, password string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf("INSERT INTO users (username, email, password) VALUES ('%s', '%s', '%s')", username, email, password))
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return err
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE TABLE `%s` ( "+
		"`uuid` char(36) NOT NULL, "+
		"`title` varchar(255) NOT NULL, "+
		"`publisher` varchar(255) DEFAULT NULL, "+
		"`category` varchar(255) DEFAULT NULL, "+
		"`author` varchar(255) DEFAULT NULL, "+
		"`page` smallint(5) UNSIGNED DEFAULT NULL, "+
		"`language` varchar(255) DEFAULT NULL, "+
		"`publication_year` year(4) DEFAULT NULL, "+
		"`isbn` varchar(20) NOT NULL, "+
		"`file_content` blob DEFAULT NULL, "+
		"`created_at` timestamp NOT NULL DEFAULT current_timestamp(), "+
		"`updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()) "+
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;", username))
	if err != nil {
		fmt.Println("Error creating user table:", err)
		return err
	}

	return nil
}
