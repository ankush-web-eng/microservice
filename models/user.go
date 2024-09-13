package models

import (
	"time"
)

type User struct {
	// gorm.Model
	ID         string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name       string     `gorm:"size:255;not null"`
	Email      string     `gorm:"size:255;not null;unique"`
	Password   string     `gorm:"not null"`
	IsVerified bool       `gorm:"default:false"`
	VerifyCode string     `gorm:"size:255;not null"`
	APIKey     *string    `gorm:"size:255"`
	Requests   int        `gorm:"default:0"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	Cloudinary Cloudinary `gorm:"foreignKey:UserID"`
	Mail       Mail       `gorm:"foreignKey:UserID"`
}

type Cloudinary struct {
	// gorm.Model
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CloudName string    `gorm:"size:255;not null"`
	APIKey    string    `gorm:"size:255;not null"`
	APISecret string    `gorm:"size:255;not null"`
	Requests  int       `gorm:"default:0"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UserID    string    `gorm:"unique;not null"` // Foreign key
}

type Mail struct {
	// gorm.Model
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email     string    `gorm:"size:255;not null"`
	Password  string    `gorm:"size:255;not null"`
	Requests  int       `gorm:"default:0"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UserID    string    `gorm:"unique;not null"` // Foreign key
}
