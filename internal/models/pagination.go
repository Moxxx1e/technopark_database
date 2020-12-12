package models

type Pagination struct {
	Since string `query:"since"`
	Limit int    `query:"limit"`
	Desc  bool   `query:"desc"`
}
