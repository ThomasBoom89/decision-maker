package database

import "gorm.io/gorm"

type Parameter struct {
	gorm.Model
	ConfigurationID uint             `gorm:"uniqueIndex:,composite:configuration_name"`
	Name            string           `gorm:"uniqueIndex:,composite:configuration_name"`
	Type            string           `gorm:"not null"`
	Comparer        string           `gorm:"not null"`
	ParameterValues []ParameterValue `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
