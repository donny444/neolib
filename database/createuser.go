package database

import (
	"context"
	"fmt"
)

func CreateUser(ctx context.Context, username string, email string, password string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return err
	}

	// Insert user into the users table
	_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO users (username, email, password) VALUES ('%s', '%s', '%s')", username, email, password))
	if err != nil {
		fmt.Println("Error inserting user:", err)
		tx.Rollback()
		return err
	}

	// Create a table for the user
	_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE TABLE `%s` ( "+
		"`title` varchar(255) NOT NULL, "+
		"`publisher` varchar(255) DEFAULT NULL, "+
		"`category` varchar(255) DEFAULT NULL, "+
		"`author` varchar(255) DEFAULT NULL, "+
		"`pages` smallint(5) UNSIGNED DEFAULT NULL, "+
		"`language` varchar(255) DEFAULT NULL, "+
		"`publication_year` year(4) DEFAULT NULL, "+
		"`isbn` varchar(20) NOT NULL, "+
		"`file` blob DEFAULT NULL, "+
		"`created_at` timestamp NOT NULL DEFAULT current_timestamp(), "+
		"`updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()) "+
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;", username))
	if err != nil {
		fmt.Println("Error creating user table:", err)
		tx.Rollback()
		return err
	}

	// Create a view for the user
	_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE VIEW `%s_view` AS SELECT * FROM `%s`", username, username))
	if err != nil {
		fmt.Println("Error creating user view:", err)
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		fmt.Println("Error committing transaction:", err)
		return err
	}

	return nil
}
