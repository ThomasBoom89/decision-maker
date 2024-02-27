package rendering

import (
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/rendering/views"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetUpRoutes(router fiber.Router, configurationRepository *database.ConfigurationRepository, databaseConnection *gorm.DB) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, Twitch!",
		})
	})

	productGroup := router.Group("/product")
	product := views.NewProduct(productGroup, databaseConnection)
	product.SetUpRoutes()

	configurationGroup := router.Group("/configuration")
	configuration := views.NewConfiguration(configurationGroup, configurationRepository, databaseConnection)
	configuration.SetUpRoutes()
}
