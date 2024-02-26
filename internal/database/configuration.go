package database

import (
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"gorm.io/gorm"
)

type Configuration struct {
	gorm.Model
	Version    uint        `gorm:"unique autoIncrement"`
	Active     bool        `gorm:"default: false"`
	Valid      bool        `gorm:"default: false"`
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
	err := result.Error
	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &configuration, nil
}

func (R *ConfigurationRepository) GetAll() ([]Configuration, error) {
	var configurations []Configuration
	err := R.database.Debug().Model(Configuration{}).Order("version DESC").Find(&configurations).Error
	if err != nil {
		return nil, err
	}

	return configurations, nil
}

func (R *ConfigurationRepository) UpdateStatus(configuration *Configuration) (*Configuration, error) {
	configuration.Active = !configuration.Active

	result := R.database.Debug().Save(configuration)
	if result.Error != nil {
		return nil, result.Error
	}
	return configuration, nil
}

func (R *ConfigurationRepository) AppendParameter(configuration *Configuration, name, parameterType, comparerType string) (*Configuration, error) {
	var dummyConfiguration Configuration
	parameter := Parameter{
		Name:     name,
		Type:     parameterType,
		Comparer: decision.Compare(comparerType),
	}
	dummyConfiguration.ID = configuration.ID
	err := R.database.Debug().Model(&dummyConfiguration).Association("Parameters").Append(&parameter)
	if err != nil {
		panic(err)
		return nil, err
	}
	configuration.Parameters = append(configuration.Parameters, parameter)

	return configuration, nil
}

func (R *ConfigurationRepository) Create(version uint) (*Configuration, error) {
	configuration := &Configuration{
		Version:    version,
		Active:     false,
		Valid:      false,
		Parameters: nil,
		Products:   nil,
	}
	err := R.database.Debug().Create(configuration).Error
	if err != nil {
		panic(err)
		return nil, err
	}

	return configuration, nil
}

func (R *ConfigurationRepository) GetNextVersion() uint {
	var result uint
	err := R.database.Debug().Model(Configuration{}).Select("MAX(version)+1 AS maxversion").Scan(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}
