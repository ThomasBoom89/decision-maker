package configuration

import (
	"errors"
	"fmt"
	"os"
)

type Loader struct {
}

func (L *Loader) LoadDatabaseConfiguration() *Database {
	return &Database{
		Host:     L.loadFromEnvironment("DATABASE_HOST"),
		Name:     L.loadFromEnvironment("DATABASE_NAME"),
		Username: L.loadFromEnvironment("DATABASE_USERNAME"),
		Password: L.loadFromEnvironment("DATABASE_PASSWORD"),
		Timezone: L.loadFromEnvironment("DATABASE_TIMEZONE"),
	}
}

func (L *Loader) loadFromEnvironment(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(errors.New(fmt.Sprint("Missing key: ", key, " from environment")))
	}

	return value
}
