package models

type Related struct {
	User   bool `query:"user,omitempty"`
	Forum  bool `query:"forum,omitempty"`
	Thread bool `query:"thread,omitempty"`
}
