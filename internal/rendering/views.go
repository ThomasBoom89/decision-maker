package rendering

import (
	"fmt"
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

	productGroup := router.Group("/product")

	productGroup.Get("/overview/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		productRepository := database.NewProductRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		products, _ := productRepository.GetByConfiguration(configuration.ID)

		fmt.Println(products)

		return ctx.Render("product/overview", fiber.Map{"Products": products})
	})

	productGroup.Get("/new/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))

		return ctx.Render("product/new", fiber.Map{
			"Parameter": configuration.Parameters,
			"Version":   version,
		})
	})

	productGroup.Post("/save", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.FormValue("version"))
		//name := ctx.FormValue("name")

		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, err := configurationRepository.GetByVersion(uint(version))
		if err != nil {
			panic(err)

		}
		parameterMap := make(map[uint]string)
		//var values []decision.ValueTypeComparer
		values := make(map[uint]decision.ValueTypeComparer)
		for _, parameter := range configuration.Parameters {
			parameterId := strconv.Itoa(int(parameter.ID))
			parameterValue := ctx.FormValue("parameter" + parameterId)
			parameterMap[parameter.ID] = parameterValue
			values[parameter.ID] = decision.ValueTypeComparer{
				Name:     parameter.Name,
				Value:    parameterValue,
				Type:     parameter.Type,
				Comparer: decision.Compare(parameter.Comparer),
			}
			//values = append(values, valueTypeComparer)
		}

		testConfigurator := decision.NewTestConfigurator()
		testConfiguration := testConfigurator.Create(values)

		//decisionMaker := decision.NewMakerForTestConfiguration()
		foobar(databaseConnection, configuration.ID, values, testConfiguration)
		// todo: test configuration logic
		//fmt.Println("testcnfig:", testConfiguration)
		//fmt.Println("parammap: ", parameterMap)
		//fmt.Println("values:", values)
		//productRepository := database.NewProductRepository(databaseConnection)
		//productRepository.InsertOne(configuration.ID, name, parameterMap, testConfiguration)

		return nil
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

/*
Check if test configuration of new product will collide with existing product
*/
func foobar(databaseConnection *gorm.DB, configurationId uint, parametersMap map[uint]decision.ValueTypeComparer, testConfiguration map[string]string) {
	productRepository := database.NewProductRepository(databaseConnection)
	products, _ := productRepository.GetByConfiguration(configurationId)

	decisionMaker := decision.NewMakerForTestConfiguration()

	for _, product := range products {
		if len(product.ParameterValues) == 0 {
			continue
		}
		fmt.Println(product.ParameterValues[0].Value)
		for _, parameterValue := range product.ParameterValues {
			fmt.Println("parameterValue:", parameterValue.Value)
			compareType := parametersMap[parameterValue.ParameterID]
			//fmt.Println("compareType:", compareType)
			fmt.Println("compareValue:", testConfiguration[compareType.Name])
			fmt.Println("valueType:", compareType.Type)
			fmt.Println("compareType:", compareType.Comparer)

			// todo: save parameter to map
			result := decisionMaker.Decide(parameterValue.Value, testConfiguration[compareType.Name], compareType.Comparer, parametersMap[parameterValue.ParameterID].Type)
			fmt.Println("result: ", result)
			// todo: find unique match and return

		}
	}
}

/*
Check if existing test configurations will match new product
*/
func barfoo() {

}
