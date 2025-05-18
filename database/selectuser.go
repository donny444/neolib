package database

import (
	"context"
	"database/sql"
)

func SelectUser(ctx context.Context, username string, email string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT username, email, password FROM users WHERE username = ? OR email = ?", username, email), nil
}
