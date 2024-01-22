package main

import (
	"errors"
	"fmt"
	"github.com/ThomasBoom89/decision-maker/internal/configuration"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func main() {
	loader := configuration.Loader{}
	databaseConfiguration := loader.LoadDatabaseConfiguration()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// middleware
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(compress.New())

	app.Static("/", "./public")

	databaseConnection, err := gorm.Open(postgres.Open(databaseConfiguration.GetPostgresDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrator := database.NewMigrator(databaseConnection)
	err = migrator.Migrate()
	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, Twitch!",
		})
	})

	app.Get("/match/:version", func(c *fiber.Ctx) error {
		version, err := strconv.Atoi(c.Params("version"))
		if err != nil {
			panic(err)
		}
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configurationByVersion, err := configurationRepository.GetByVersion(uint(version))
		if err != nil {
			//return http.StatusBadRequest
			return err
		}

		queryMap := make(map[string]string)
		for _, parameter := range configurationByVersion.Parameters {
			value := c.Query(strings.ToLower(parameter.Name), "foobar")
			if value == "foobar" {
				//return http.StatusBadRequest
				return errors.New(fmt.Sprint("query param missing: ", parameter.Name))
			}
			queryMap[parameter.Name] = value
		}

		return c.Render("testversion", fiber.Map{
			"Version": version,
			"Params":  c.Params("*"),
			"Debug":   queryMap,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
