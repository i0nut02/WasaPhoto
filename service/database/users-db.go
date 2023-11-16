package database

import (
	"encoding/json"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/sirupsen/logrus"
)

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
