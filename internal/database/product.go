package database

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ConfigurationID   uint                        `gorm:"uniqueIndex:,composite:configuration_name"`
	Name              string                      `gorm:"not null;uniqueIndex:,composite:configuration_name"`
	ParameterValues   []ParameterValue            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Changelog         []ChangelogProductParameter `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TestConfiguration TestConfiguration           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ProductRepository struct {
	database *gorm.DB
}

func NewProductRepository(database *gorm.DB) *ProductRepository {
	return &ProductRepository{database: database}
}

func (P *ProductRepository) GetByConfiguration(configurationId uint) ([]Product, error) {
	var product []Product
	err := P.database.Debug().Model(Product{}).Where("configuration_id = ?", configurationId).Preload("TestConfiguration").Preload("ParameterValues").Find(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (P *ProductRepository) InsertOne(configurationId uint, name string, parameterValuesMap map[uint]string, testConfigurationMap map[string]string) {
	var parameterValues []ParameterValue
	for key, value := range parameterValuesMap {
		parameterValues = append(parameterValues, ParameterValue{
			ParameterID: key,
			Value:       value,
		})
	}

	product := Product{
		Model:           gorm.Model{},
		ConfigurationID: configurationId,
		Name:            name,
		ParameterValues: parameterValues,
		Changelog:       nil,
		TestConfiguration: TestConfiguration{
			Configuration: testConfigurationMap,
		},
	}

	P.database.Create(&product)
}

func (P *ProductRepository) GetProductIdsByConfiguration(configurationId uint) ([]uint, error) {
	rows, _ := P.database.Debug().Table("products").Select("id").Where("configuration_id = ?", configurationId).Rows()
	var fooSlice []uint
	for rows.Next() {
		var fooBar uint
		err := rows.Scan(&fooBar)
		if err != nil {
			return nil, err
		}
		fooSlice = append(fooSlice, fooBar)
	}

	return fooSlice, nil
}

func (P *ProductRepository) Delete(id uint) error {
	err := P.database.Debug().Delete(&Product{}, id).Error
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func (P *ProductRepository) GetOne(id uint) (Product, error) {
	var product Product
	err := P.database.Debug().Preload("ParameterValues").Find(&product, id).Error
	if err != nil {
		panic(err)
		return Product{}, err
	}

	return product, nil
}
