package database

import "context"

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
