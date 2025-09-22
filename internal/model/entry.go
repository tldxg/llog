package model

import (
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Body string
}

type Summary struct {
	Date    string `json:"date"`
	Summary string `json:"summary"`
}
