package database

import (
	"context"
	"database/sql"
)

func SelectBook(ctx context.Context, uuid string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT uuid, title, publisher, category, author, page, language, publication_year, isbn FROM books WHERE uuid = ?", uuid), nil
}
