package rendering

import (
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/rendering/views"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(
	router fiber.Router,
	productRepository *database.ProductRepository,
	configurationRepository *database.ConfigurationRepository,
	testConfigurationRepository *database.TestConfigurationRepository,
) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, Twitch!",
		})
	})

	productGroup := router.Group("/product")
	product := views.NewProduct(productGroup, productRepository, configurationRepository, testConfigurationRepository)
	product.SetUpRoutes()

	configurationGroup := router.Group("/configuration")
	configuration := views.NewConfiguration(configurationGroup, configurationRepository)
	configuration.SetUpRoutes()
}
