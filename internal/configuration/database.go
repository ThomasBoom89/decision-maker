package configuration

import "fmt"

const PostgresDsnFormat = "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=%s"

type Database struct {
	Host     string
	Name     string
	Username string
	Password string
	Timezone string
}

func (D *Database) GetPostgresDSN() string {
	return fmt.Sprintf(PostgresDsnFormat, D.Host, D.Username, D.Password, D.Name, D.Timezone)
}
