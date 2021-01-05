package models

type Pagination struct {
	Limit int    `query:"limit"`
	Desc  bool   `query:"desc"`
}
