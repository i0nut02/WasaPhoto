package database

import (
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/google/uuid"
)

func (db *appdbimpl) DoLogin(username string) (json string, err error) {
	var count int

	err = db.c.QueryRow(`SELECT Count() FROM users WHERE username = ?`, username).Scan(&count)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error checking the number of user with input username in the db, datails: %w", err)
	}

	if count == 0 {
		newUUID := uuid.New()

		userId := newUUID.String()

		_, err = db.c.Exec(`INSERT INTO users (id, username) VALUES (?, ?)`, userId, username)

		if err != nil {
			return components.InternalServerError, fmt.Errorf("error inserting on the database a new user, details: %w", err)
		}
	}

	user_json := components.User{Username: username}
	data, err := user_json.ToJson()

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error converting user to Json, details: %w", err)
	}

	return string(data), nil
}