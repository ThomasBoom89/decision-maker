package main

import (
	"errors"
	"fmt"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// middleware
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(compress.New())

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, Twitch!",
		})
	})

	dsn := "host=decision-maker-database user=root password=root dbname=decision-maker port=5432 sslmode=disable TimeZone=Europe/Berlin"
	//dsn := "host=localhost user=root password=root dbname=decision-maker port=5433 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrator := database.NewMigrator(db)
	err = migrator.Migrate()
	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":3002"))
}

func getFromEnv(key string) string {
	value, err := os.LookupEnv(key)
	if err {
		panic(errors.New(fmt.Sprint("Missing key: ", key, " from environment")))
	}

	return value
}
