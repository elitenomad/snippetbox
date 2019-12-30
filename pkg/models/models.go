package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrInvalidCredentials = errors.New("models: Invalid credentials")
var ErrDuplicateEmail = errors.New("models: Duplicate Email")

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}


type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
	Active bool
}