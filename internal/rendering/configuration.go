package rendering

import (
	"fmt"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"github.com/ThomasBoom89/decision-maker/internal/rendering/views"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"strconv"
)

type Configuration struct {
	router                  fiber.Router
	configurationRepository *database.ConfigurationRepository
	configurationView       *views.Configuration
}

func NewConfiguration(router fiber.Router, configurationRepository *database.ConfigurationRepository, configurationView *views.Configuration) *Configuration {
	return &Configuration{
		router:                  router,
		configurationRepository: configurationRepository,
		configurationView:       configurationView,
	}
}

func (C *Configuration) SetUpRoutes() {

	C.router.Get("/new", func(ctx *fiber.Ctx) error {
		configuration := database.Configuration{}
		configuration.ID = 1
		editConfiguration := C.configurationView.Edit(C.getParameterTypes(), C.getCompareTypes(), configuration)

		return adaptor.HTTPHandler(templ.Handler(editConfiguration))(ctx)
	})

	C.router.Delete("/:id", func(ctx *fiber.Ctx) error {
		id, _ := strconv.Atoi(ctx.Params("id"))
		err := C.configurationRepository.Delete(uint(id))
		if err != nil {
			return err
		}

		return nil
	})

	C.router.Get("/copy/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configuration, _ := C.configurationRepository.GetByVersion(uint(version))

		nextVersion := C.configurationRepository.GetNextVersion()
		newConfiguration, _ := C.configurationRepository.Create(nextVersion)
		for _, parameter := range configuration.Parameters {
			newConfiguration, _ = C.configurationRepository.AppendParameter(newConfiguration, parameter.Name, parameter.Type, string(parameter.Comparer))
		}

		return ctx.Redirect(fmt.Sprintf("/configuration/edit/%d", newConfiguration.Version))
	})

	C.router.Get("/show/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configuration, _ := C.configurationRepository.GetByVersion(uint(version))
		showConfiguration := C.configurationView.Show(*configuration)

		return adaptor.HTTPHandler(templ.Handler(showConfiguration))(ctx)
	})

	C.router.Get("/overview", func(ctx *fiber.Ctx) error {
		configurations, err := C.configurationRepository.GetAll()
		if err != nil {
			return err
		}

		overview := C.configurationView.Overview(configurations)

		return adaptor.HTTPHandler(templ.Handler(overview))(ctx)
	})

	C.router.Get("/edit/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configuration, _ := C.configurationRepository.GetByVersion(uint(version))
		if configuration.Active {
			return ctx.Redirect("/configuration/overview", 302)
		}
		editConfiguration := C.configurationView.Edit(C.getParameterTypes(), C.getCompareTypes(), *configuration)

		return adaptor.HTTPHandler(templ.Handler(editConfiguration))(ctx)
	})

	C.router.Get("/status/change/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configuration, _ := C.configurationRepository.GetByVersion(uint(version))
		configuration, _ = C.configurationRepository.UpdateStatus(configuration)
		overviewTableRow := C.configurationView.OverviewTableRow(*configuration)

		return adaptor.HTTPHandler(templ.Handler(overviewTableRow))(ctx)
	})

	C.router.Get("/comparer", func(ctx *fiber.Ctx) error {
		parameterType := ctx.Query("type")
		compareTypes := decision.GetCompareTypes()
		compareSelect := C.configurationView.GetCompareTypeSelect(compareTypes[parameterType])

		return adaptor.HTTPHandler(templ.Handler(compareSelect))(ctx)
	})

	C.router.Post("/create/parameter/:version?", func(ctx *fiber.Ctx) error {
		parameterType := ctx.FormValue("type")
		comparerType := ctx.FormValue("comparer")
		name := ctx.FormValue("name")
		var configuration *database.Configuration
		if ctx.Params("version") == "" {
			version := C.configurationRepository.GetNextVersion()
			configuration, _ = C.configurationRepository.Create(version)
		} else {
			version, _ := strconv.Atoi(ctx.Params("version"))
			configuration, _ = C.configurationRepository.GetByVersion(uint(version))
		}
		configuration, _ = C.configurationRepository.AppendParameter(configuration, name, parameterType, comparerType)
		editConfiguration := C.configurationView.EditForm(C.getParameterTypes(), C.getCompareTypes(), *configuration)

		return adaptor.HTTPHandler(templ.Handler(editConfiguration))(ctx)
	})
}

func (C *Configuration) getCompareTypes() []decision.Compare {
	return []decision.Compare{
		decision.GreaterThan,
		decision.GreaterEqual,
		decision.LowerThan,
		decision.LowerEqual,
		decision.Equal,
		decision.NotEqual,
		decision.Range,
	}
}

func (C *Configuration) getParameterTypes() []string {
	return []string{
		decision.Integer,
		decision.String,
		decision.Float,
		decision.DateTime,
		decision.Boolean,
		decision.Time,
		decision.Date,
	}
}
