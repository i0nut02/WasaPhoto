package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func (db *appdbimpl) GetUsernameFromId(id string) (username string, err error) {
	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf(components.BadRequestError)
		}
		return "", err
	}
	return username, nil
}

func (db *appdbimpl) IsValid(userId string, username string) (is_valid bool, err error) {
	var count int

	err = db.c.QueryRow(`SELECT Count() FROM users WHERE id = ? and username = ?`, userId, username).Scan(&count)

	is_valid = count == 1

	if err != nil {
		return false, fmt.Errorf("error during the query, details: %w", err)
	}
	return is_valid, nil
}

func (db *appdbimpl) SearchUserBySubString(searcher string, searched_string string) (maches string, err error) {
	res, err := db.c.Query(`SELECT username FROM users WHERE username LIKE '%'||?||'%' EXCEPT SELECT banisher FROM bans WHERE banished = ? EXCEPT SELECT ?;`, searched_string, searcher, searcher)

	defer func() {
		if res != nil {
			err := res.Close()

			if err != nil {
				logrus.Errorf("error closing query result, details: %v", err)
			}
		}
	}()

	if err != nil {
		return components.InternalServerError, err
	}

	var users []components.User
	for res.Next() {
		if err = res.Err(); err != nil {
			return components.InternalServerError, err
		}

		var user components.User

		err = res.Scan(&user.Username)

		if err != nil {
			return components.InternalServerError, err
		}

		users = append(users, user)
	}

	data, err := json.Marshal(users)

	if err != nil {
		return components.InternalServerError, err
	}
	if string(data) == "null" {
		return "[]", nil
	}
	return string(data), nil
}

func (db *appdbimpl) SetUsername(id string, old_username string, new_username string) (data string, err error) {
	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE id = ? and username = ?`, new_username, id, old_username)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.Code == sqlite3.ErrConstraint {
				return components.ConflictError, fmt.Errorf("already exist a user with the new username")
			}
		}
		return components.InternalServerError, fmt.Errorf("error checking if the new username is taken")
	}
	var user components.User
	user.Username = new_username

	res, err := user.ToJson()

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error parsing the new username")
	}

	return string(res), nil
}
