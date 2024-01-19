package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/components"
)

func (db *appdbimpl) BanUser(idBanisher string, idBanished string) (data string, err error) {
	_, err = db.c.Exec(`INSERT OR IGNORE INTO bans (banisher, banished) VALUES (?, ?)`, idBanisher, idBanished)

	if err != nil {
		return components.InternalServerError, err
	}

	var user components.User

	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, idBanished).Scan(&user.Username)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := user.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}

	return string(jsonData), err
}

func (db *appdbimpl) UnbanUser(idBanisher string, idBanished string) (data string, err error) {
	_, err = db.c.Exec(`DELETE FROM bans WHERE banisher = ? AND banished = ?`, idBanisher, idBanished)

	if err != nil {
		return components.InternalServerError, err
	}

	var user components.User

	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, idBanished).Scan(&user.Username)

	if err != nil {
		return components.InternalServerError, err
	}

	jsonData, err := user.ToJson()

	if err != nil {
		return components.InternalServerError, err
	}

	return string(jsonData), err
}
