package views

import (
	"fmt"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type Configuration struct {
	router             fiber.Router
	databaseConnection *gorm.DB
}

func NewConfiguration(router fiber.Router, databaseConnection *gorm.DB) *Configuration {
	return &Configuration{
		router:             router,
		databaseConnection: databaseConnection,
	}
}

func (C *Configuration) SetUpRoutes() {

	C.router.Get("/new", func(ctx *fiber.Ctx) error {

		return ctx.Render("configuration/new", fiber.Map{
			"Title":          "New Configuration",
			"ParameterTypes": C.getParameterTypes(),
			"CompareTypes":   C.getCompareTypes(),
		})
	})

	C.router.Get("/copy/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(C.databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))

		nextVersion := configurationRepository.GetNextVersion()
		newConfiguration, _ := configurationRepository.Create(nextVersion)
		for _, parameter := range configuration.Parameters {
			newConfiguration, _ = configurationRepository.AppendParameter(newConfiguration, parameter.Name, parameter.Type, string(parameter.Comparer))
		}

		return ctx.Redirect(fmt.Sprintf("/configuration/edit/%d", newConfiguration.Version))
	})

	C.router.Get("/show/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(C.databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))

		return ctx.Render("configuration/show", fiber.Map{
			"Title":         "Show",
			"Configuration": configuration,
		})
	})

	C.router.Get("/overview", func(ctx *fiber.Ctx) error {
		configurationRepository := database.NewConfigurationRepository(C.databaseConnection)
		configurations, err := configurationRepository.GetAll()
		if err != nil {
			return err
		}

		return ctx.Render("configuration/overview", fiber.Map{
			"Configurations": configurations,
		})
	})

	C.router.Get("/edit/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(C.databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		if configuration.Active {
			return ctx.Redirect("/configuration/overview", 302)
		}

		return ctx.Render("configuration/edit", fiber.Map{
			"Title":          "Edit configuration",
			"ParameterTypes": C.getParameterTypes(),
			"CompareTypes":   C.getCompareTypes(),
			"Configuration":  configuration,
		})
	})

	C.router.Get("/status/change/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(C.databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		configuration, _ = configurationRepository.UpdateStatus(configuration)
		return ctx.Render("configuration/overview_row", fiber.Map{
			"Version": configuration.Version,
			"Active":  configuration.Active,
		})
	})

	C.router.Get("/comparer", func(ctx *fiber.Ctx) error {
		parameterType := ctx.Query("type")
		compareTypes := decision.GetCompareTypes()
		return ctx.Render("configuration/parameter_compare", fiber.Map{
			"CompareTypes": compareTypes[parameterType],
		})
	})

	C.router.Post("/create/parameter/:version?", func(ctx *fiber.Ctx) error {
		parameterType := ctx.FormValue("type")
		comparerType := ctx.FormValue("comparer")
		name := ctx.FormValue("name")
		configurationRepository := database.NewConfigurationRepository(C.databaseConnection)
		var configuration *database.Configuration
		if ctx.Params("version") == "" {
			version := configurationRepository.GetNextVersion()
			configuration, _ = configurationRepository.Create(version)
		} else {
			version, _ := strconv.Atoi(ctx.Params("version"))
			configuration, _ = configurationRepository.GetByVersion(uint(version))
		}
		configuration, _ = configurationRepository.AppendParameter(configuration, name, parameterType, comparerType)

		return ctx.Render("configuration/edit_form", fiber.Map{
			"Title":          "Edit configuration",
			"ParameterTypes": C.getParameterTypes(),
			"CompareTypes":   C.getCompareTypes(),
			"Configuration":  configuration,
		})
	})
}

func (C *Configuration) getCompareTypes() []string {
	return []string{"gt", "ge", "lt", "le", "eq", "ne"}
}

func (C *Configuration) getParameterTypes() []string {
	return []string{"int", "string", "float", "datetime", "bool"}
}
