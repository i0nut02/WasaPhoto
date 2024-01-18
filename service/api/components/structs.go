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
	NumFollowed  int    `json:"num_following"`
	Following    bool   `json:"following"`
	IsBanished    bool   `json:"is_banished"`
}

func (u *UserProfile) ToJson() ([]byte, error) {
	return json.MarshalIndent(u, "", " ")
}

type Post struct {
	UploadTime  string `json:"upload_time"`
	Author      string `json:"author"`
	NumLikes    int    `json:"num_likes"`
	NumComments int    `json:"num_comments"`
	LikedPhoto  bool   `json:"liked_photo"`
	PhotoID     string `json:"photo_id"`
	PhotoFile   string `json:"photo_file"`
	Description string `json:"description"`
}

func (p *Post) ToJson() ([]byte, error) {
	return json.MarshalIndent(p, "", " ")
}

type Photo struct {
	File        string `json:"file"`
	Description string `json:"description"`
}

func (p *Photo) ToJson() ([]byte, error) {
	return json.MarshalIndent(p, "", " ")
}

type Comment struct {
	Id         string `json:"id"`
	Author     string `json:"author"`
	Content    string `json:"content"`
	UploadTime string `json:"upload_time"`
}

func (c *Comment) ToJson() ([]byte, error) {
	return json.MarshalIndent(c, "", " ")
}

type CommentContent struct {
	Content string `json:"comment"`
}
