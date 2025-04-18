package database

import (
	"context"
	"database/sql"
)

func FindUser(ctx context.Context, usernameOrEmail string) (*sql.Row, error) {
	return db.QueryRowContext(ctx, "SELECT username, email, password FROM users WHERE username = ? OR email = ?", usernameOrEmail, usernameOrEmail), nil
}
