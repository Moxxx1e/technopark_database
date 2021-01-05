package models

type Vote struct {
	ThreadID uint64 `json:"thread_id"`
	UserID   uint64 `json:"user_id"`
	Likes    bool   `json:"likes"`
}
