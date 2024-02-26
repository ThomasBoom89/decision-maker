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
		testConfigurationOldProducts := foobar(databaseConnection, configuration.ID, values, testConfiguration)
		//barfoo(databaseConnection, configuration)
		oldTestConfigurationNewProduct := foobarfoo(databaseConnection, configuration, values)
		// todo: test configuration logic
		//fmt.Println("testcnfig:", testConfiguration)
		//fmt.Println("parammap: ", parameterMap)
		//fmt.Println("values:", values)
		//productRepository := database.NewProductRepository(databaseConnection)
		//productRepository.InsertOne(configuration.ID, name, parameterMap, testConfiguration)

		return ctx.Render("product/diff", fiber.Map{
			"TestConfiguration":              testConfiguration,
			"TestConfigurationOldProducts":   testConfigurationOldProducts,
			"OldTestConfigurationNewProduct": oldTestConfigurationNewProduct,
		})
	})

	configurationGroup := router.Group("/configuration")

	configurationGroup.Get("/new", func(ctx *fiber.Ctx) error {

		return ctx.Render("configuration/new", fiber.Map{
			"Title":          "New Configuration",
			"ParameterTypes": getParameterTypes(),
			"CompareTypes":   getCompareTypes(),
		})
	})

	configurationGroup.Get("/copy/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))

		nextVersion := configurationRepository.GetNextVersion()
		newConfiguration, _ := configurationRepository.Create(nextVersion)
		for _, parameter := range configuration.Parameters {
			newConfiguration, _ = configurationRepository.AppendParameter(newConfiguration, parameter.Name, parameter.Type, string(parameter.Comparer))
		}

		return ctx.Render("configuration/copy", fiber.Map{
			"Title":          "Copy",
			"Configuration":  newConfiguration,
			"ParameterTypes": getParameterTypes(),
			"CompareTypes":   getCompareTypes(),
		})
	})

	configurationGroup.Get("/show/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))

		return ctx.Render("configuration/show", fiber.Map{
			"Title":         "Show",
			"Configuration": configuration,
		})
	})

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

		return ctx.Render("configuration/edit", fiber.Map{
			"Title":          "Edit configuration",
			"ParameterTypes": getParameterTypes(),
			"CompareTypes":   getCompareTypes(),
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

	configurationGroup.Post("/create/parameter/:version?", func(ctx *fiber.Ctx) error {
		parameterType := ctx.FormValue("type")
		comparerType := ctx.FormValue("comparer")
		name := ctx.FormValue("name")
		configurationRepository := database.NewConfigurationRepository(databaseConnection)
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
			"ParameterTypes": getParameterTypes(),
			"CompareTypes":   getCompareTypes(),
			"Configuration":  configuration,
		})
	})
}

func getCompareTypes() []string {
	return []string{"gt", "ge", "lt", "le", "eq", "ne"}
}

func getParameterTypes() []string {
	return []string{"int", "string", "float", "datetime", "bool"}
}

/*
Check if test configuration of new product will collide with existing product
*/
func foobar(databaseConnection *gorm.DB, configurationId uint, parametersMap map[uint]decision.ValueTypeComparer, testConfiguration map[string]string) map[string][]Result {
	productRepository := database.NewProductRepository(databaseConnection)
	products, _ := productRepository.GetByConfiguration(configurationId)
	fmt.Println(testConfiguration)
	decisionMaker := decision.NewMakerForTestConfiguration()

	result := make(map[string][]Result)
	for _, product := range products {
		if len(product.ParameterValues) == 0 {
			continue
		}
		fmt.Println(product.ParameterValues[0].Value)
		for _, parameterValue := range product.ParameterValues {
			compareType := parametersMap[parameterValue.ParameterID]
			decisionResult := decisionMaker.Decide(parameterValue.Value, testConfiguration[compareType.Name], compareType.Comparer, parametersMap[parameterValue.ParameterID].Type)
			result[product.Name] = append(result[product.Name], Result{
				ParameterName: compareType.Name,
				TestValue:     testConfiguration[compareType.Name],
				ProductValue:  parameterValue.Value,
				CompareType:   compareType.Type,
				Result:        decisionResult,
			})
			// todo: find unique match and return
		}
	}

	return result
}

/*
*
Check if existing test configurations will match new product
*/
func foobarfoo(databaseConnection *gorm.DB, configuration *database.Configuration, comparerMap map[uint]decision.ValueTypeComparer) map[string][]Result {
	productRepository := database.NewProductRepository(databaseConnection)
	productIds, _ := productRepository.GetProductIdsByConfiguration(configuration.ID)
	products, _ := productRepository.GetByConfiguration(configuration.ID)
	testConfigurationRepository := database.NewTestConfigurationRepository(databaseConnection)
	testConfigurations := testConfigurationRepository.GetByProductIds(productIds)
	decisionMaker := decision.NewMakerForTestConfiguration()
	productsMap := make(map[uint]string)
	for _, product := range products {
		productsMap[product.ID] = product.Name
	}
	result := make(map[string][]Result)
	for _, testConfiguration := range testConfigurations {
		for _, parameter := range configuration.Parameters {
			comparer := comparerMap[parameter.ID]
			decisionResult := decisionMaker.Decide(comparer.Value, testConfiguration.Configuration[comparer.Name], comparer.Comparer, comparer.Type)
			result[productsMap[testConfiguration.ProductID]] = append(result[productsMap[testConfiguration.ProductID]], Result{
				ParameterName: parameter.Name,
				TestValue:     testConfiguration.Configuration[comparer.Name],
				ProductValue:  comparer.Value,
				CompareType:   comparer.Type,
				Result:        decisionResult,
			})
		}
	}

	return result
}

/*
Check if existing test configurations will match existing products
*/
func barfoo(databaseConnection *gorm.DB, configuration *database.Configuration) {
	productRepository := database.NewProductRepository(databaseConnection)
	productIds, _ := productRepository.GetProductIdsByConfiguration(configuration.ID)
	products, _ := productRepository.GetByConfiguration(configuration.ID)
	testConfigurationRepository := database.NewTestConfigurationRepository(databaseConnection)
	testConfigurations := testConfigurationRepository.GetByProductIds(productIds)
	decisionMaker := decision.NewMakerForTestConfiguration()

	productsMap := make(map[uint]map[uint]decision.ValueTypeComparer)
	for _, product := range products {
		parameterMap := make(map[uint]database.ParameterValue)
		for _, parameterValue := range product.ParameterValues {
			parameterMap[parameterValue.ParameterID] = parameterValue
		}
		productMap := make(map[uint]decision.ValueTypeComparer)
		for _, parameter := range configuration.Parameters {
			productMap[parameter.ID] = decision.ValueTypeComparer{
				Name:     parameter.Name,
				Value:    parameterMap[parameter.ID].Value,
				Type:     parameter.Type,
				Comparer: parameter.Comparer,
			}
		}

		productsMap[product.ID] = productMap
	}
	for _, testConfiguration := range testConfigurations {
		for productId, product := range productsMap {
			if productId == testConfiguration.ProductID {
				continue
			}
			for _, comparer := range product {
				fmt.Println("product parameter value:", comparer.Value)
				//fmt.Println("compareType:", compareType)
				fmt.Println("compareValue:", testConfiguration.Configuration[comparer.Name])
				fmt.Println("valueType:", comparer.Type)
				fmt.Println("compareType:", comparer.Comparer)
				result := decisionMaker.Decide(comparer.Value, testConfiguration.Configuration[comparer.Name], comparer.Comparer, comparer.Type)
				fmt.Println("result: ", result)
				//foobar(databaseConnection, configuration.ID, product, testConfiguration.Configuration)
			}
		}
	}

}

type Result struct {
	ParameterName string
	TestValue     string
	ProductValue  string
	CompareType   string
	Result        bool
}
