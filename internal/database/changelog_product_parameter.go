package database

import "gorm.io/gorm"

type ChangelogProductParameter struct {
	gorm.Model
	ProductID uint
	Name      string `gorm:"not null"`
	OldValue  string `gorm:"not null"`
	NewValue  string `gorm:"not null"`
}
