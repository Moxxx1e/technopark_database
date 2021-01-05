package models

import "time"

type Post struct {
	ID       uint64    `json:"id"`
	Parent   uint64    `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   uint64    `json:"thread"`
	Created  time.Time `json:"created"`
}
