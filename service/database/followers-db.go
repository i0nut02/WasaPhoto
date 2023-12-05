package database

import (
	"encoding/json"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) GetFollowers(username string, requestedId string) (data string, err error) {
	rows, err := db.c.Query(`SELECT U1.username
							 FROM followers F INNER JOIN users U1 ON F.follower = U1.id INNER JOIN users U2 ON F.followed = U2.id
							 WHERE U2.username = ?
							 EXCEPT
							 SELECT U.username
							 FROM bans B INNER JOIN users U ON B.banisher = U.id
							 WHERE B.banished = ?`, username, requestedId)
	if err != nil {
		return components.InternalServerError, err
	}

	defer rows.Close()

	var followerList []components.User

	for rows.Next() {
		if err = rows.Err(); err != nil {
			return components.InternalServerError, err
		}
		var follower components.User

		err = rows.Scan(&follower.Username)

		if err != nil {
			return components.InternalServerError, err
		}

		followerList = append(followerList, follower)
	}

	jsonData, err := json.Marshal(followerList)

	if err != nil {
		return components.InternalServerError, err
	}

	if string(jsonData) == EmptyJsonArray {
		return "[]", nil
	}
	return string(jsonData), nil
}

func (db *appdbimpl) GetFollowing(username string, requestedId string) (data string, err error) {
	rows, err := db.c.Query(`SELECT U2.username
							 FROM followers F INNER JOIN users U1 ON F.follower = U1.id INNER JOIN users U2 ON F.followed = U2.id
							 WHERE U1.username = ?
							 EXCEPT
							 SELECT U.username
							 FROM bans B INNER JOIN users U ON B.banisher = U.id
							 WHERE B.banished = ?`, username, requestedId)
	if err != nil {
		return components.InternalServerError, err
	}

	defer rows.Close()

	var followerList []components.User

	for rows.Next() {
		if err = rows.Err(); err != nil {
			return components.InternalServerError, err
		}
		var follower components.User

		err = rows.Scan(&follower.Username)

		if err != nil {
			return components.InternalServerError, err
		}

		followerList = append(followerList, follower)
	}

	jsonData, err := json.Marshal(followerList)

	if err != nil {
		return components.InternalServerError, err
	}

	if string(jsonData) == EmptyJsonArray {
		return "[]", nil
	}
	return string(jsonData), nil
}

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
