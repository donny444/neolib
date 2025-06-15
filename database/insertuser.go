package database

import (
	"context"
	"fmt"
	"log"
)

func InsertUser(ctx context.Context, username string, email string, password string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Error starting transaction: ", err)
	}

	// Insert user into the users table
	_, err = tx.ExecContext(ctx, "INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error inserting user: ", err)
	}

	// Create a table for the user
	_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE TABLE `%s` ( "+
		"`title` VARCHAR(255) NOT NULL, "+
		"`publisher` VARCHAR(255) DEFAULT NULL, "+
		"`category` VARCHAR(255) DEFAULT NULL, "+
		"`author` VARCHAR(255) DEFAULT NULL, "+
		"`pages` SMALLINT(5) UNSIGNED DEFAULT NULL, "+
		"`language` VARCHAR(255) DEFAULT NULL, "+
		"`publication_year` YEAR(4) DEFAULT NULL, "+
		"`is_read` BOOL DEFAULT 0, "+
		"`isbn` VARCHAR(20) NOT NULL, "+
		"`file` BLOB DEFAULT NULL, "+
		"`path` VARCHAR(32) DEFAULT NULL, "+
		"`created_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(), "+
		"`updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()) "+
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;", username))
	if err != nil {
		tx.Rollback()
		log.Fatal("Error creating user table: ", err)
	}

	// Create a view for the user
	_, err = tx.ExecContext(ctx, fmt.Sprintf("CREATE VIEW `%s_view` AS SELECT * FROM `%s`", username, username))
	if err != nil {
		tx.Rollback()
		log.Fatal("Error creating user view: ", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Fatal("Error committing transaction: ", err)
	}

	return nil
}
