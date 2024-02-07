package rendering

import (
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func SetUpRoutes(router fiber.Router, databaseConnection *gorm.DB) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, Twitch!",
		})
	})

	configurationGroup := router.Group("/configuration")

	configurationGroup.Get("/overview", func(ctx *fiber.Ctx) error {
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configurations, err := configurationRepository.GetAll()
		if err != nil {
			return err
		}

		return ctx.Render("configuration/overview", fiber.Map{
			"Configurations": configurations,
		})
	})

	configurationGroup.Get("/edit/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		if configuration.Active {
			return ctx.Redirect("/configuration/overview", 302)
		}
		// todo: refactor
		parameterTypes := []string{"int", "string", "float", "datetime", "bool"}
		compareTypes := []string{"gt", "ge", "lt", "le", "eq", "ne"}
		return ctx.Render("configuration/edit", fiber.Map{
			"Title":          "Edit configuration",
			"ParameterTypes": parameterTypes,
			"CompareTypes":   compareTypes,
			"Configuration":  configuration,
		})
	})

	configurationGroup.Get("/status/change/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		configuration, _ = configurationRepository.UpdateStatus(configuration)
		return ctx.Render("configuration/overview_row", fiber.Map{
			"Version": configuration.Version,
			"Active":  configuration.Active,
		})
	})

	configurationGroup.Get("/comparer", func(ctx *fiber.Ctx) error {
		parameterType := ctx.Query("type")
		compareTypes := decision.GetCompareTypes()
		return ctx.Render("configuration/parameter_compare", fiber.Map{
			"CompareTypes": compareTypes[parameterType],
		})
	})

	configurationGroup.Post("/create/parameter/:version", func(ctx *fiber.Ctx) error {
		parameterType := ctx.FormValue("type")
		comparerType := ctx.FormValue("comparer")
		name := ctx.FormValue("name")
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		configuration, _ = configurationRepository.AppendParameter(configuration, name, parameterType, comparerType)
		// todo: refactor
		parameterTypes := []string{"int", "string", "float", "datetime", "bool"}
		compareTypes := []string{"gt", "ge", "lt", "le", "eq", "ne"}
		return ctx.Render("configuration/edit_form", fiber.Map{
			"Title":          "Edit configuration",
			"ParameterTypes": parameterTypes,
			"CompareTypes":   compareTypes,
			"Configuration":  configuration,
		})
	})
}
