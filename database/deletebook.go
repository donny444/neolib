package database

import "context"

func DeleteBook(ctx context.Context, isbn string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM books WHERE isbn = ?", isbn)
	return err
}
