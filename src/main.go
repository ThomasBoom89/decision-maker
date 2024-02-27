package main

import (
	"github.com/ThomasBoom89/decision-maker/internal/configuration"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/rendering"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	loader := configuration.Loader{}
	databaseConfiguration := loader.LoadDatabaseConfiguration()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:             engine,
		EnablePrintRoutes: true,
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

	configurationRepository := database.NewConfigurationRepository(databaseConnection)
	productRepository := database.NewProductRepository(databaseConnection)
	testConfigurationRepository := database.NewTestConfigurationRepository(databaseConnection)
	// rendering frontend views
	views := app.Group("/")
	rendering.SetUpRoutes(views, productRepository, configurationRepository, testConfigurationRepository)

	// api
	//api := app.Group("/api")
	log.Fatal(app.Listen(":3000"))
}
