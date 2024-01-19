package database

import "gorm.io/gorm"

type TestConfiguration struct {
	gorm.Model
	ProductID     uint
	Configuration map[string]string `gorm:"not null;serializer:json"`
}
