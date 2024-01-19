package database

import "gorm.io/gorm"

type Configuration struct {
	gorm.Model
	Version    uint        `gorm:"unique autoIncrement"`
	Parameters []Parameter `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Products   []Product   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ConfigurationRepository struct {
}
