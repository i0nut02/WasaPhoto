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

func (db *appdbimpl) GetIdFromUsername(username string) (id string, err error) {
	err = db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf(components.BadRequestError)
		}
		return "", err
	}
	return id, nil
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

func (db *appdbimpl) ValidUsername(username string) (is_valid bool, err error) {
	err = db.c.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)`, username).Scan(&is_valid)

	if err != nil {
		return false, err
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

func (db *appdbimpl) IsBanished(banisher string, banished string) (answer bool, err error) {
	var banisherId string
	var banishedId string

	err = db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, banisher).Scan(&banisherId)
	if err != nil {
		return false, err
	}

	err = db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, banished).Scan(&banishedId)
	if err != nil {
		return false, err
	}

	err = db.c.QueryRow(`SELECT EXISTS (SELECT 1 FROM bans WHERE banisher = ? and banished = ?)`, banisher, banished).Scan(&answer)

	if err != nil {
		return false, err
	}
	return answer, nil
}

func (db *appdbimpl) TakeProfile(username string, usernameProfile string) (profile string, err error) {
	var userProfile components.UserProfile
	var usernameId string
	var usernameProfileId string

	err = db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&usernameId)

	if err != nil {
		return components.InternalServerError, err
	}

	err = db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, usernameProfile).Scan(&usernameProfileId)

	if err != nil {
		return components.InternalServerError, err
	}

	err = db.c.QueryRow(`
		SELECT 
			X1.username,
			X2.num_photos,
			X3.num_followers,
			X4.num_followed,
			(SELECT EXISTS (SELECT 1 FROM followers WHERE followed = ? AND follower = ?)) AS is_following
		FROM 
			(SELECT username FROM users WHERE id = ?) X1,
			(SELECT COUNT(*) as num_photos FROM posts WHERE author_id = ?) X2,
			(SELECT COUNT(*) as num_followers FROM followers WHERE followed = ?) X3,
			(SELECT COUNT(*) as num_followed FROM followers WHERE follower = ?) X4
			`, usernameProfileId, usernameId, usernameProfileId, usernameProfileId, usernameProfileId, usernameProfileId).Scan(
		&userProfile.Username,
		&userProfile.NumPhotos,
		&userProfile.NumFollowers,
		&userProfile.NumFollowed,
		&userProfile.Following,
	)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := userProfile.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}

	return string(jsonData), nil
}

func (db *appdbimpl) GetStream(id string, offset int, limit int) (data string, err error) {
	rows, err := db.c.Query(`
		SELECT 
			P.upload_time AS upload_time, 
			U1.username AS author, 
			COUNT(L.liker) AS num_likes, 
			COUNT(C.comment_id) AS num_comments, 
			EXISTS (
				SELECT 1 
				FROM likes 
				WHERE post_id = P.id AND liker = ?
			) AS liked_photo, 
			P.id AS photo_id, 
			P.file AS photo_file, 
			P.description AS description
		FROM (
			SELECT followed AS user 
			FROM followers 
			WHERE follower = ?
			EXCEPT
			SELECT banner 
			FROM bans 
			WHERE banned = ?
		) U
		INNER JOIN posts P ON P.author_id = U.user
		INNER JOIN users U1 ON U1.id = U.user
		INNER JOIN likes L ON L.post_id = P.id
		INNER JOIN comments C ON C.post_id = P.id
		GROUP BY P.id
		ORDER BY P.upload_time DESC
		OFFSET ?
		LIMIT ?
		`, id, id, id, offset, limit)

	if err != nil {
		return components.InternalServerError, err
	}

	defer rows.Close()

	var posts []components.Post

	for rows.Next() {
		if err = rows.Err(); err != nil {
			return components.InternalServerError, err
		}

		var post components.Post
		err = rows.Scan(
			&post.UploadTime,
			&post.Author,
			&post.NumLikes,
			&post.NumComments,
			&post.LikedPhoto,
			&post.PhotoID,
			&post.PhotoFile,
			&post.Description,
		)

		if err != nil {
			return components.InternalServerError, err
		}

		posts = append(posts, post)
	}

	res, err := json.Marshal(posts)

	if err != nil {
		return components.InternalServerError, err
	}
	if string(res) == "null" {
		return "[]", nil
	}
	return string(res), nil
}
