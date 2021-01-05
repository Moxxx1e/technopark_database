package models

import "time"

type Thread struct {
	ID      uint64    `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created"`
}
