package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var ErrorNoRecord = errors.New("no matching records found")

type Snippets struct {
	ID      int
	Uuid    string
	Title   string
	Content string
	Created time.Time `gorm:"autoCreateTime"`
	Expires time.Time
}

type Users struct {
	ID       int
	Uuid     string
	Name     string
	Email    string
	Password string
	Created  time.Time `gorm:"autoCreateTime"`
}

func (user *Users) BeforeCreate(tx *gorm.DB) error {
	user.Uuid = uuid.New().String()
	return nil
}

func (snippet *Snippets) BeforeCreate(tx *gorm.DB) error {
	snippet.Uuid = uuid.New().String()
	return nil
}
