package model

import (
	"errors"
)

// Various errors the models might return.
var (
	ErrNotAllowed = errors.New("request not allowed")
	// ErrLoginWrongPass wrong password
	ErrLoginWrongPass = errors.New("Wrong Password")
)
