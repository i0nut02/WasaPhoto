package components

import (
	"encoding/json"
	"regexp"
)

type User struct {
	Username string `json:"username_string"`
}

func (u *User) ToJson() ([]byte, error) {
	return json.MarshalIndent(u, "", " ")
}

func (u *User) IsValidUsername() (match bool, err error) {
	pattern := "^[a-zA-Z0-9_.]{3,20}$"

	regex, err := regexp.Compile(pattern)

	if err != nil {
		return false, err
	}
	return regex.MatchString(u.Username), nil
}

type UserId struct {
	UserId string `json:"user_id"`
}

func (u *UserId) ToJson() ([]byte, error) {
	return json.MarshalIndent(u, "", " ")
}

type UserProfile struct {
	Username     string `json:"username"`
	NumPhotos    int    `json:"num_photos"`
	NumFollowers int    `json:"num_followers"`
	NumFollowed  int    `json:"num_followed"`
	Following    bool   `json:"following"`
}

func (u *UserProfile) ToJson() ([]byte, error) {
	return json.MarshalIndent(u, "", " ")
}
