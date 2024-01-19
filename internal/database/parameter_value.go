package database

import "gorm.io/gorm"

type ParameterValue struct {
	gorm.Model
	ParameterID uint
	ProductID   uint
	Value       string `gorm:"not null"`
}
