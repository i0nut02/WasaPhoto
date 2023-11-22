package database

import (
	"encoding/json"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) ValidPostAuthor(postId string, author string) (isValid bool, err error) {
	err = db.c.QueryRow(`SELECT (EXISTS (SELECT 1 
										FROM posts P 
										INNER JOIN users U 
										ON P.author_id = U.id 
										WHERE P.id = ? AND
											U.username = ?))`, postId, author).Scan(&isValid)
	if err != nil {
		return false, err
	}
	return isValid, nil
}

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
							P.photo AS photo_file, 
							P.description AS description
						FROM users U
						INNER JOIN posts P ON P.author_id = U.id
						LEFT JOIN likes L ON L.post_id = P.id
						LEFT JOIN comments C ON C.post_id = P.id
						WHERE P.author_id = ?
						GROUP BY P.id
						ORDER BY P.upload_time DESC
						`, requesterId, id)

	if err != nil {
		return components.InternalServerError, err
	}
	defer rows.Close()

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
	_, err = db.c.Exec(`DELETE FROM posts WHERE id = ?`, post_id)

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
						P.photo AS photo_file, 
						P.description AS description
					FROM users U
					INNER JOIN posts P ON P.author_id = U.id
					LEFT JOIN likes L ON L.post_id = P.id
					LEFT JOIN comments C ON C.post_id = P.id
					WHERE P.id = ?
					GROUP BY P.id;
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

func (db *appdbimpl) GetLikes(postId string, userId string) (data string, err error) {
	rows, err := db.c.Query(`SELECT U.username 
							 FROM users U
							 INNER JOIN likes L
							 ON U.id = L.liker
							 WHERE L.post_id = ?
							 EXCEPT
							 SELECT U.username
							 FROM users U
							 INNER JOIN bans B
							 ON U.id = B.banisher
							 WHERE B.banished = ?`, postId, userId)

	if err != nil {
		return components.InternalServerError, err
	}

	var userList []components.User

	for rows.Next() {
		if err = rows.Err(); err != nil {
			return components.InternalServerError, err
		}
		var user components.User

		rows.Scan(&user.Username)

		userList = append(userList, user)
	}

	jsonData, err := json.Marshal(userList)

	if err != nil {
		return components.InternalServerError, err
	}

	if string(jsonData) == "null" {
		return "[]", nil
	}
	return string(jsonData), nil
}

func (db *appdbimpl) LikePhoto(postId string, likerId string) (data string, err error) {
	_, err = db.c.Exec(`INSERT INTO likes (post_id, liker) VALUES (?, ?)`, postId, likerId)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); !ok || !(sqliteErr.Code == sqlite3.ErrConstraint) {
			return components.InternalServerError, err
		}
	}

	var user components.User

	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, likerId).Scan(&user.Username)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := user.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}
	return string(jsonData), nil
}

func (db *appdbimpl) UnlikePhoto(postId string, likerId string) (data string, err error) {
	_, err = db.c.Exec(`DELETE FROM likes WHERE post_id = ? AND liker = ?`, postId, likerId)

	if err != nil {
		return components.InternalServerError, err
	}

	var user components.User

	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, likerId).Scan(&user)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := user.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}
	return string(jsonData), nil
}

func (db *appdbimpl) GetComments(postId string, userId string) (data string, err error) {
	rows, err := db.c.Query(` 
							SELECT
								C.comment_id,
								U1.username,
								C.content,
								C.upload_time
							FROM
								comments C
								INNER JOIN 
								(SELECT id AS id FROM users EXCEPT SELECT banisher FROM bans WHERE banished = ?) U
								ON U.id = C.user_id
								INNER JOIN users U1 ON U1.id = U.id
							WHERE
								C.post_id = ?
							`, userId, postId)
	if err != nil {
		return components.InternalServerError, err
	}

	var comments []components.Comment
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return components.InternalServerError, err
		}

		var comment components.Comment

		err = rows.Scan(&comment.Id,
			&comment.Author,
			&comment.Content,
			&comment.UploadTime)

		if err != nil {
			return components.InternalServerError, err
		}

		comments = append(comments, comment)
	}

	res, err := json.Marshal(comments)

	if err != nil {
		return components.InternalServerError, err
	}

	if string(res) == "null" {
		return "[]", nil
	}

	return string(res), nil
}

func (db *appdbimpl) CommentPhoto(postId string, userId string, comment string) (data string, err error) {
	newUUID := uuid.New()
	timestampFormatted := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	var commentComp components.Comment

	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, userId).Scan(&commentComp.Author)

	if err != nil {
		return components.InternalServerError, err
	}

	_, err = db.c.Exec(`INSERT INTO comments (comment_id, post_id, user_id, content, upload_time) VALUES (?, ?, ?, ?, ?)`, newUUID, postId, userId, comment, timestampFormatted)

	if err != nil {
		return components.InternalServerError, err
	}

	commentComp.Content = comment
	commentComp.Id = newUUID.String()

	jsonData, err := commentComp.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}

	return string(jsonData), nil
}

func (db *appdbimpl) UncommentPhoto(commentId string) (data string, err error) {
	_, err = db.c.Exec(`DELETE FROM comments WHERE comment_id = ?`, commentId)

	if err != nil {
		return components.InternalServerError, err
	}

	return "", nil
}
