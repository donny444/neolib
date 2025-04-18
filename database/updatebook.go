package database

import "context"

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
