package database

import "gorm.io/gorm"

type TestConfiguration struct {
	gorm.Model
	ProductID     uint
	Configuration map[string]string `gorm:"not null;serializer:json"`
}

type TestConfigurationRepository struct {
	Database *gorm.DB
}

func NewTestConfigurationRepository(database *gorm.DB) *TestConfigurationRepository {
	return &TestConfigurationRepository{
		Database: database,
	}
}

func (T *TestConfigurationRepository) InsertOne(productId uint, configuration map[string]string) *TestConfiguration {

	testConfiguration := &TestConfiguration{
		Model:         gorm.Model{},
		ProductID:     productId,
		Configuration: configuration,
	}

	T.Database.Debug().Create(testConfiguration)

	return testConfiguration
}

func (T *TestConfigurationRepository) GetByProductIds(productIds []uint) []TestConfiguration {
	var testConfigurations []TestConfiguration

	T.Database.Debug().Where("product_id IN ?", productIds).Find(&testConfigurations)

	return testConfigurations
}
