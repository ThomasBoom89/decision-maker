package database

import "gorm.io/gorm"

type Migrator struct {
	database *gorm.DB
}

func NewMigrator(database *gorm.DB) *Migrator {
	return &Migrator{
		database: database,
	}
}

func (m *Migrator) Migrate() error {

	err := m.database.Debug().AutoMigrate(
		&Configuration{},
		&Parameter{},
		&Product{},
		&ParameterValue{},
		&ChangelogProductParameter{},
		&TestConfiguration{},
	)
	if err != nil {
		return err
	}

	return nil
}
