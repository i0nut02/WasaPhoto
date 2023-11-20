package database

import (
	"encoding/json"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/google/uuid"
)

func (db *appdbimpl) UploadPhoto(photo *components.Photo, id string) (data string, err error) {
	timestampFormatted := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	newUUID := uuid.New()

	_, err = db.c.Exec(`INSERT INTO posts (id, photo, author_id, upload_time, description)
					 VALUES (?, ?, ?, ?, ?)`, newUUID, photo.File, id, timestampFormatted, photo.Description)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := photo.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}
	return string(jsonData), err
}

func (db *appdbimpl) GetUserPosts(id string, requesterId string) (data string, err error) {
	var posts []components.Post

	rows, err := db.c.Query(`
						SELECT 
							P.upload_time AS upload_time, 
							U.username AS author, 
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
						FROM users U
						INNER JOIN posts P ON P.author_id = U.id
						INNER JOIN likes L ON L.post_id = P.id
						INNER JOIN comments C ON C.post_id = P.id
						WHERE P.author = ?
						ORDER BY P.upload_time DESC
						`, requesterId, id)
	if err != nil {
		return components.InternalServerError, err
	}

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

func (db *appdbimpl) DeletePhoto(post_id string) (data string, err error) {
	_, err = db.c.Exec(`DELETE FROM posts WHERE post_id = ?`, post_id)

	if err != nil {
		return components.InternalServerError, err
	}

	return components.NoContent, nil
}

func (db *appdbimpl) GetPost(post_id string, id string) (data string, err error) {
	var post components.Post
	err = db.c.QueryRow(`
						SELECT 
							P.upload_time AS upload_time, 
							U.username AS author, 
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
						FROM users U
						INNER JOIN posts P ON P.author_id = U.id
						INNER JOIN likes L ON L.post_id = P.id
						INNER JOIN comments C ON C.post_id = P.id
						WHERE P.post_id = ?
						`, id, post_id).Scan(
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

	jsonData, err := post.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}

	return string(jsonData), nil
}
