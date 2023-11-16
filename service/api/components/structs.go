package components

import (
	"encoding/json"
)

type User struct {
	Username string `json:"username_string"`
}

func (u *User) ToJson() ([]byte, error) {
	return json.MarshalIndent(u, "", " ")
}
