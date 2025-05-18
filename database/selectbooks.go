package database

import (
	"context"
	"database/sql"
	"fmt"
)

func SelectBooks(ctx context.Context, category *string, username string) (*sql.Rows, error) {
	if category != nil {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT isbn, title FROM `%s_view` WHERE category = ?", username), *category)
	} else {
		return db.QueryContext(ctx, fmt.Sprintf("SELECT isbn, title FROM `%s_view`", username))
	}
}
