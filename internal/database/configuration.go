package database

import (
	"errors"
	"fmt"
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
	affectedRows := R.database.Debug().Model(Configuration{}).Where("version = ?", version).Preload("Parameters").First(&configuration).RowsAffected
	if affectedRows == 0 {
		return nil, errors.New(fmt.Sprint("no configuration with version ", version, " found!"))
	}

	return &configuration, nil
}
