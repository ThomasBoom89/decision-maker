package views

import (
	"fmt"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type Product struct {
	router             fiber.Router
	databaseConnection *gorm.DB
}

func NewProduct(router fiber.Router, databaseConnection *gorm.DB) *Product {
	return &Product{
		router:             router,
		databaseConnection: databaseConnection,
	}
}

func (P *Product) SetUpRoutes() {

	P.router.Get("/overview/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(P.databaseConnection)
		productRepository := database.NewProductRepository(P.databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))
		products, _ := productRepository.GetByConfiguration(configuration.ID)

		fmt.Println(products)

		return ctx.Render("product/overview", fiber.Map{"Products": products})
	})

	P.router.Get("/new/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configurationRepository := database.NewConfigurationRepository(P.databaseConnection)
		configuration, _ := configurationRepository.GetByVersion(uint(version))

		return ctx.Render("product/new", fiber.Map{
			"Parameter": configuration.Parameters,
			"Version":   version,
		})
	})

	P.router.Post("/save", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.FormValue("version"))
		//name := ctx.FormValue("name")

		configurationRepository := database.NewConfigurationRepository(P.databaseConnection)
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
		testConfigurationOldProducts := P.foobar(configuration.ID, values, testConfiguration)
		//barfoo(databaseConnection, configuration)
		oldTestConfigurationNewProduct := P.foobarfoo(configuration, values)
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
}

/*
Check if test configuration of new product will collide with existing product
*/
func (P *Product) foobar(configurationId uint, parametersMap map[uint]decision.ValueTypeComparer, testConfiguration map[string]string) map[string][]Result {
	productRepository := database.NewProductRepository(P.databaseConnection)
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
func (P *Product) foobarfoo(configuration *database.Configuration, comparerMap map[uint]decision.ValueTypeComparer) map[string][]Result {
	productRepository := database.NewProductRepository(P.databaseConnection)
	productIds, _ := productRepository.GetProductIdsByConfiguration(configuration.ID)
	products, _ := productRepository.GetByConfiguration(configuration.ID)
	testConfigurationRepository := database.NewTestConfigurationRepository(P.databaseConnection)
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
func (P *Product) barfoo(configuration *database.Configuration) {
	productRepository := database.NewProductRepository(P.databaseConnection)
	productIds, _ := productRepository.GetProductIdsByConfiguration(configuration.ID)
	products, _ := productRepository.GetByConfiguration(configuration.ID)
	testConfigurationRepository := database.NewTestConfigurationRepository(P.databaseConnection)
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
