package database

import (
	"context"
	"database/sql"
)

func SelectBooks(ctx context.Context, category *string, username string) (*sql.Rows, error) {
	// username := r.Context()
	// username := ctx.Value()
	if category != nil {
		return db.QueryContext(ctx, "SELECT isbn, title FROM ? WHERE category = ?", username, *category)
	} else {
		return db.QueryContext(ctx, "SELECT isbn, title FROM ?", username)
	}
}
