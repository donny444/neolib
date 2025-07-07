package database

import (
	"context"
	"fmt"
)

func InsertBook(
	ctx context.Context,
	username string,
	isbn string,
	title string,
	publisher *string,
	category *string,
	author *string,
	pages *string,
	language *string,
	publicationYear *string,
	fileContent []byte,
	fileExtension *string,
) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO `%s` (isbn, title, publisher, category, author, pages, language, publication_year, file_content, file_extension)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		username),
		isbn,
		title,
		publisher,
		category,
		author,
		pages,
		language,
		publicationYear,
		fileContent,
		fileExtension,
	)
	return err
}
