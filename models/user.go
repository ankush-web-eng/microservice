package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"uniqueIndex"`
	Password   string
	IsActive   bool
	activeCode string
}