package database

import (
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"gorm.io/gorm"
)

type Parameter struct {
	gorm.Model
	ConfigurationID uint             `gorm:"uniqueIndex:,composite:configuration_name"`
	Name            string           `gorm:"uniqueIndex:,composite:configuration_name"`
	Type            string           `gorm:"not null"`
	Comparer        decision.Compare `gorm:"not null"`
	ParameterValues []ParameterValue `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ParameterRepository struct {
	database *gorm.DB
}

func NewParameterRepository(database *gorm.DB) *ParameterRepository {
	return &ParameterRepository{database: database}
}

func (P *ParameterRepository) GetParameterByConfigurationId(id uint) ([]Parameter, error) {
	var parameters []Parameter
	err := P.database.Debug().Model(Parameter{}).Where("configuration_id = ?", id).Find(parameters).Error
	if err != nil {
		return nil, err
	}

	return parameters, nil
}
