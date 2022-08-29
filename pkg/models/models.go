package models

import (
	"errors"
	"time"
)

var ErrorNoRecord = errors.New("no matching records found")

type Snippets struct {
	ID      int
	Title   string
	Content string
	Created time.Time `gorm:"autoCreateTime"`
	Expires time.Time
}
