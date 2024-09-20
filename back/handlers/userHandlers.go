package handlers

import (
	"context"

	"github.com/code-troopers/postitsonline/database"
)

func CreateUser(user database.User) error {
	rows, err := database.DB.Query(context.Background(), "SELECT id FROM users WHERE id = $1", user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		if id != "" {
			return nil
		}
	}
	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO users (id, given_name, family_name, email, picture) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.GivenName, user.FamilyName, user.Email, user.Picture)
	if err != nil {
		return err
	}
	return nil
}
