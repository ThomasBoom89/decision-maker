package database

import (
	"gorm.io/gorm"
)

type Configuration struct {
	gorm.Model
	Version    uint        `gorm:"unique autoIncrement"`
	Parameters []Parameter `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Products   []Product   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ConfigurationRepository struct {
	database *gorm.DB
}

func NewConfigurationRepository(database *gorm.DB) *ConfigurationRepository {
	return &ConfigurationRepository{database: database}
}

func (R *ConfigurationRepository) GetByVersion(version uint) (*Configuration, error) {
	var configuration Configuration
	result := R.database.Debug().Model(Configuration{}).Where("version = ?", version).Preload("Parameters").First(&configuration)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &configuration, nil
}
