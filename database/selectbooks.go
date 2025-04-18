package database

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectBooks(ctx context.Context, category *string) (*sql.Rows, error) {
	if category != nil {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT uuid, title, isbn FROM books WHERE category = '%s'", *category))
	} else {
		return db.QueryContext(ctx, "SELECT uuid, title, isbn FROM books")
	}
}
