package database

import (
	"context"
	"database/sql"
)

func SelectBooks(ctx context.Context, category *string) (*sql.Rows, error) {
	if category != nil {
		return db.QueryContext(ctx, "SELECT isbn, title FROM books WHERE category = ?", *category)
	} else {
		return db.QueryContext(ctx, "SELECT isbn, title FROM books")
	}
}
