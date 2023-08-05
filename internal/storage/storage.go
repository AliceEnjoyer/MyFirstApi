package storage

import "errors"

// 12) тут вся логика как в телеграм боте

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type storage interface {
}
