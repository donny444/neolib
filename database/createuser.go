package database

import (
	"context"
	"fmt"
)

func CreateUser(ctx context.Context, username string, email string, password string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf("INSERT INTO users (username, email, password) VALUES ('%s', '%s', '%s')", username, email, password))
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return err
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE TABLE `%s` ( "+
		"`uuid` char(36) NOT NULL, "+
		"`title` varchar(255) NOT NULL, "+
		"`publisher` varchar(255) DEFAULT NULL, "+
		"`category` varchar(255) DEFAULT NULL, "+
		"`author` varchar(255) DEFAULT NULL, "+
		"`page` smallint(5) UNSIGNED DEFAULT NULL, "+
		"`language` varchar(255) DEFAULT NULL, "+
		"`publication_year` year(4) DEFAULT NULL, "+
		"`isbn` varchar(20) NOT NULL, "+
		"`file_content` blob DEFAULT NULL, "+
		"`created_at` timestamp NOT NULL DEFAULT current_timestamp(), "+
		"`updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()) "+
		"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;", username))
	if err != nil {
		fmt.Println("Error creating user table:", err)
		return err
	}

	return nil
}
