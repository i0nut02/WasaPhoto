package database

import (
	"encoding/json"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) FollowUser(followerId string, followedId string) (data string, err error) {
	var user components.User

	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, followedId).Scan(&user.Username)

	if err != nil {
		return components.InternalServerError, err
	}

	_, err = db.c.Exec(`INSERT OR IGNORE INTO followers (followed, follower) VALUES (?, ?)`, followedId, followerId)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := user.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}

	return string(jsonData), nil
}

func (db *appdbimpl) UnfollowUser(followerId string, followedId string) (data string, err error) {
	_, err = db.c.Exec(`DELETE FROM followers WHERE followed = ? AND follower = ?`, followedId, followerId)

	if err != nil {
		return components.InternalServerError, err
	}
	return "", nil
}
