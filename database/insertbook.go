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
	path *string,
) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO `%s` (isbn, title, publisher, category, author, pages, language, publication_year, file, path)"+
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
		path,
	)
	return err
}
