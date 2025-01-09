package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupDatabase() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/neolib")
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
	return db.QueryRowContext(ctx, "SELECT title, publisher, category, author, page, language, publication_year, isbn FROM books WHERE uuid = ?", uuid), nil
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
