/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	DoLogin(username string) (string, error)

	IsValid(userId string, username string) (bool, error)

	SearchUserBySubString(searcher string, serch_string string) (matches string, err error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// i shoud also set PRAGMA to True (?)
		usersDatabase := `CREATE TABLE IF NOT EXISTS users (
							id TEXT NOT NULL PRIMARY KEY, 
							username TEXT NOT NULL UNIQUE
							);`
		followersDatabase := `CREATE TABLE IF NOT EXISTS followers (
								followed TEXT NOT NULL,
								follower TEXT NOT NULL,
								FOREIGN KEY (followed) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
								FOREIGN KEY (follower) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
								);`
		bansDataBase := `CREATE TABLE IF NOT EXISTS bans (
							banisher TEXT NOT NULL,
							banished TEXT NOT NULL,
							FOREIGN KEY (banisher) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
							FOREIGN KEY (banished) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
							);`
		postsDatabase := `CREATE TABLE IF NOT EXISTS posts (
							id TEXT NOT NULL PRIMARY KEY,
							photo BLOB NOT NULL,
							author_id TEXT NOT NULL,
							upload_time TEXT NOT NULL,
							description TEXT,
							FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
							);`
		likesDatabase := `CREATE TABLE IF NOT EXISTS likes (
							post_id TEXT NOT NULL,
							liker TEXT NOT NULL,
							FOREIGN KEY (liker) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
							FOREIGN KEY (liker) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
							)`
		commentsDatabase := `CREATE TABLE IF NOT EXISTS comments (
								comment_id TEXT NOT NULL PRIMARY KEY,
								post_id TEXT NOT NULL,
								user_id TEXT NOT NULL,
								content TEXT NOT NULL,
								FOREIGN KEY (post_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
								FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
		)`
		_, err = db.Exec(usersDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(followersDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(bansDataBase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(postsDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(likesDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
		_, err = db.Exec(commentsDatabase)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
