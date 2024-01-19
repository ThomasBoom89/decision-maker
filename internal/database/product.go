package database

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ConfigurationID   uint                        `gorm:"uniqueIndex:,composite:configuration_name"`
	Name              string                      `gorm:"not null;uniqueIndex:,composite:configuration_name"`
	ParameterValues   []ParameterValue            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Changelog         []ChangelogProductParameter `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TestConfiguration TestConfiguration           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
