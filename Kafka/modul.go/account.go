package modul

import "time"

type Account struct {
	Id        int        `json:"id"`
	Balance   float64    `json:"balance"`
	Currency  string     `json:"currency"`
	IsLocked  bool       `json:"is_locked"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
