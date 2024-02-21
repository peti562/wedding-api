package model

import (
	"time"
)

type Invite struct {
	Id            string     `db:"id" json:"id"`
	Name          string     `db:"name" json:"name"`
	Greeting      string     `db:"greeting" json:"greeting"`
	Language      string     `db:"lang" json:"lang"`
	MaxAdults     int        `db:"max_adults" json:"max_adults"`
	MaxChildren   int        `db:"max_children" json:"max_children"`
	RSVP          bool       `db:"rsvp" json:"rsvp"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	FirstOpenedAt *time.Time `db:"first_opened_at" json:"first_opened_at"`
	LastOpenedAt  *time.Time `db:"last_opened_at" json:"last_opened_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}
type Request struct {
	InviteId string     `db:"invite_id" json:"invite_id"`
	RSVP     bool       `db:"rsvp" json:"rsvp"`
	Data     []Attendee `json:"data"`
}
type Attendee struct {
	Id        string     `db:"id" json:"id"`
	InviteId  string     `db:"invite_id" json:"invite_id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Phone     string     `db:"phone" json:"phone"`
	Age       int        `db:"age" json:"age"`
	Active    bool       `db:"active" json:"active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
