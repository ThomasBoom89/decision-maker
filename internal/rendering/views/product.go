package views

import (
	"fmt"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	router                      fiber.Router
	productRepository           *database.ProductRepository
	configurationRepository     *database.ConfigurationRepository
	testConfigurationRepository *database.TestConfigurationRepository
}

func NewProduct(
	router fiber.Router,
	productRepository *database.ProductRepository,
	configurationRepository *database.ConfigurationRepository,
	testConfigurationRepository *database.TestConfigurationRepository,
) *Product {
	return &Product{
		router:                      router,
		productRepository:           productRepository,
		configurationRepository:     configurationRepository,
		testConfigurationRepository: testConfigurationRepository,
	}
}

func (P *Product) SetUpRoutes() {

	P.router.Delete("/:id", func(ctx *fiber.Ctx) error {
		id, _ := strconv.Atoi(ctx.Params("id"))

		return P.productRepository.Delete(uint(id))
	})

	P.router.Get("/overview/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configuration, _ := P.configurationRepository.GetByVersion(uint(version))
		products, _ := P.productRepository.GetByConfiguration(configuration.ID)

		return ctx.Render("product/overview", fiber.Map{
			"Version":  version,
			"Products": products,
		})
	})

	P.router.Get("/edit/:id", func(ctx *fiber.Ctx) error {
		id, _ := strconv.Atoi(ctx.Params("id"))
		product, _ := P.productRepository.GetOne(uint(id))
		configuration, _ := P.configurationRepository.GetById(product.ConfigurationID)
		parameterValues := make(map[uint]database.ParameterValue)
		for _, parameterValue := range product.ParameterValues {
			parameterValues[parameterValue.ParameterID] = parameterValue
		}
		parameters := configuration.Parameters
		for key, parameter := range parameters {
			if value, ok := parameterValues[parameter.ID]; ok {
				if parameter.Comparer == decision.Range {
					splittedValue := strings.Split(value.Value, decision.RangeSeparator)
					minValue := value
					maxValue := value
					minValue.Value = splittedValue[0]
					maxValue.Value = splittedValue[1]
					minValue.Value = P.getParameterFormat(splittedValue[0], parameter.Type)
					maxValue.Value = P.getParameterFormat(splittedValue[1], parameter.Type)
					parameters[key].ParameterValues = []database.ParameterValue{minValue, maxValue}

					continue
				}
				value.Value = P.getParameterFormat(value.Value, parameter.Type)
				parameters[key].ParameterValues = []database.ParameterValue{value}
			}
		}

		return ctx.Render("product/new", fiber.Map{
			"Parameter": parameters,
			"Version":   configuration.Version,
			"Name":      product.Name,
			"ID":        product.ID,
		})
	})

	P.router.Get("/new/:version", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.Params("version"))
		configuration, _ := P.configurationRepository.GetByVersion(uint(version))

		return ctx.Render("product/new", fiber.Map{
			"Parameter": configuration.Parameters,
			"Version":   version,
		})
	})

	P.router.Post("/save/:id?", func(ctx *fiber.Ctx) error {
		version, _ := strconv.Atoi(ctx.FormValue("version"))
		id, _ := strconv.Atoi(ctx.Params("id", "0"))
		name := ctx.FormValue("name")
		configuration, err := P.configurationRepository.GetByVersion(uint(version))
		if err != nil {
			panic(err)

		}
		var products []database.Product
		if id > 0 {
			products, err = P.productRepository.GetByConfigurationExceptProduct(configuration.ID, uint(id))
		} else {
			products, err = P.productRepository.GetByConfiguration(configuration.ID)
		}
		parameterMap := make(map[uint]string)
		values := make(map[uint]decision.ValueTypeComparer)
		for _, parameter := range configuration.Parameters {
			parameterId := strconv.Itoa(int(parameter.ID))
			parameterValue := ctx.FormValue("parameter" + parameterId)
			parameterValue = P.getParsedParameter(parameterValue, parameter.Type)
			if parameter.Comparer == decision.Range {
				parameterValue += decision.RangeSeparator + P.getParsedParameter(ctx.FormValue("range"+parameterId), parameter.Type)
			}
			parameterMap[parameter.ID] = parameterValue
			values[parameter.ID] = decision.ValueTypeComparer{
				Name:     parameter.Name,
				Value:    parameterValue,
				Type:     parameter.Type,
				Comparer: parameter.Comparer,
			}
		}

		testConfigurator := decision.NewTestConfigurator()
		testConfiguration := testConfigurator.Create(values)
		testConfigurationOldProducts := P.foobar(products, values, testConfiguration)
		oldTestConfigurationNewProduct := P.foobarfoo(products, values)

		insert := len(testConfigurationOldProducts) == 0 && len(oldTestConfigurationNewProduct) == 0
		if insert {
			if id > 0 {
				product, _ := P.productRepository.GetOne(uint(id))
				product.TestConfiguration.Configuration = testConfiguration
				for key, parameterValue := range product.ParameterValues {
					product.ParameterValues[key].Value = values[parameterValue.ParameterID].Value
				}
				_ = P.productRepository.Update(product)
			} else {
				P.productRepository.InsertOne(configuration.ID, name, parameterMap, testConfiguration)
			}
		}

		return ctx.Render("product/diff", fiber.Map{
			"TestConfiguration":              testConfiguration,
			"TestConfigurationOldProducts":   testConfigurationOldProducts,
			"OldTestConfigurationNewProduct": oldTestConfigurationNewProduct,
			"Insert":                         insert,
		})
	})
}

func (P *Product) getParameterFormat(parameterValue string, parameterType string) string {
	switch parameterType {
	case decision.DateTime:
		timestamp, err := strconv.Atoi(parameterValue)
		if err != nil {
			panic(err)
		}

		return time.Unix(int64(timestamp), 0).Format("2006-01-02T15:04")
	case decision.Date:
		timestamp, err := strconv.Atoi(parameterValue)
		if err != nil {
			panic(err)
		}

		return time.Unix(int64(timestamp), 0).Format("2006-01-02")
	case decision.Time:
		timestamp, err := strconv.Atoi(parameterValue)
		if err != nil {
			panic(err)
		}

		return time.Unix(int64(timestamp), 0).Format("15:04")
	default:
		return parameterValue
	}
}

func (P *Product) getParsedParameter(parameterValue string, parameterType string) string {
	switch parameterType {
	case decision.DateTime:
		dateTime, err := time.Parse("2006-01-02T15:04", parameterValue)
		if err != nil {
			panic(err)
		}

		return strconv.Itoa(int(dateTime.Unix()))
	case decision.Date:
		date, err := time.Parse("2006-01-02", parameterValue)
		if err != nil {
			panic(err)
		}

		return strconv.Itoa(int(date.Unix()))
	case decision.Time:
		time, err := time.Parse("15:04", parameterValue)
		if err != nil {
			panic(err)
		}

		return strconv.Itoa(int(time.Unix()))
	default:
		return parameterValue
	}
}

/*
Check if test configuration of new product will collide with existing product
*/
func (P *Product) foobar(products []database.Product, parametersMap map[uint]decision.ValueTypeComparer, testConfiguration map[string]string) map[string][]Result {
	decisionMaker := decision.NewMakerForTestConfiguration()

	result := make(map[string][]Result)
	for _, product := range products {
		if len(product.ParameterValues) == 0 {
			continue
		}
		for _, parameterValue := range product.ParameterValues {
			compareType := parametersMap[parameterValue.ParameterID]
			decisionResult := decisionMaker.Decide(parameterValue.Value, testConfiguration[compareType.Name], compareType.Comparer, parametersMap[parameterValue.ParameterID].Type)
			if decisionResult == false {
				continue
			}
			result[product.Name] = append(result[product.Name], Result{
				ParameterName: compareType.Name,
				TestValue:     testConfiguration[compareType.Name],
				ProductValue:  parameterValue.Value,
				CompareType:   compareType.Type,
				Comparer:      compareType.Comparer,
				Result:        decisionResult,
			})
			// todo: find unique match and return
		}
		if len(product.ParameterValues) != len(result[product.Name]) {
			delete(result, product.Name)
		}
	}

	return result
}

/*
*
Check if existing test configurations will match new product
*/
func (P *Product) foobarfoo(products []database.Product, comparerMap map[uint]decision.ValueTypeComparer) map[string][]Result {
	decisionMaker := decision.NewMakerForTestConfiguration()
	productsMap := make(map[uint]string)
	for _, product := range products {
		productsMap[product.ID] = product.Name
	}
	result := make(map[string][]Result)
	for _, product := range products {
		for key, parameter := range comparerMap {
			comparer := comparerMap[key]
			decisionResult := decisionMaker.Decide(comparer.Value, product.TestConfiguration.Configuration[comparer.Name], comparer.Comparer, comparer.Type)
			if decisionResult == false {
				continue
			}
			result[productsMap[product.TestConfiguration.ProductID]] = append(result[productsMap[product.TestConfiguration.ProductID]], Result{
				ParameterName: parameter.Name,
				TestValue:     product.TestConfiguration.Configuration[comparer.Name],
				ProductValue:  comparer.Value,
				CompareType:   comparer.Type,
				Comparer:      comparer.Comparer,
				Result:        decisionResult,
			})
		}
		if len(comparerMap) != len(result[productsMap[product.TestConfiguration.ProductID]]) {
			delete(result, productsMap[product.TestConfiguration.ProductID])
		}
	}

	return result
}

/*
Check if existing test configurations will match existing products
*/
func (P *Product) barfoo(configuration *database.Configuration) {
	productIds, _ := P.productRepository.GetProductIdsByConfiguration(configuration.ID)
	products, _ := P.productRepository.GetByConfiguration(configuration.ID)
	testConfigurations := P.testConfigurationRepository.GetByProductIds(productIds)
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
	Comparer      decision.Compare
	Result        bool
}
