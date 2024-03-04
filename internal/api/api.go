package api

import (
	"errors"
	"github.com/ThomasBoom89/decision-maker/internal/database"
	"github.com/ThomasBoom89/decision-maker/internal/decision"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Api struct {
	router                  fiber.Router
	configurationRepository *database.ConfigurationRepository
	productRepository       *database.ProductRepository
	decisionMaker           *decision.Maker
}

func NewApi(
	router fiber.Router,
	configurationRepository *database.ConfigurationRepository,
	productRepository *database.ProductRepository,
	decisionMaker *decision.Maker,
) *Api {
	return &Api{
		router:                  router,
		configurationRepository: configurationRepository,
		productRepository:       productRepository,
		decisionMaker:           decisionMaker,
	}
}

func (A *Api) SetUpRoutes() {
	A.router.Get("/decide/:version", func(ctx *fiber.Ctx) error {
		version, err := strconv.Atoi(ctx.Params("version"))
		if err != nil {
			return err
		}
		configuration, err := A.configurationRepository.GetByVersion(uint(version))
		if err != nil {
			return err
		}

		queryParameters := make(map[uint]decision.ValueTypeComparer)
		for _, parameter := range configuration.Parameters {
			get := ctx.Query(parameter.Name)
			if get == "" {
				return errors.New("query param " + parameter.Name + " missing!")
			}
			queryParameters[parameter.ID] = decision.ValueTypeComparer{
				Name:     parameter.Name,
				Value:    get,
				Type:     parameter.Type,
				Comparer: parameter.Comparer,
			}
		}

		products, err := A.productRepository.GetByConfiguration(configuration.ID)
		if err != nil {
			return err
		}

		for _, product := range products {
			result := A.checkProductValues(product, queryParameters)
			if result == true {
				return ctx.JSON(product)
			}
		}

		return ctx.JSON(fiber.Map{})
	})
}

func (A *Api) checkProductValues(product database.Product, queryParameters map[uint]decision.ValueTypeComparer) bool {
	for _, parameterValue := range product.ParameterValues {
		queryParameter := queryParameters[parameterValue.ParameterID]
		result := A.decisionMaker.Decide(parameterValue.Value, queryParameter.Value, queryParameter.Comparer, queryParameter.Type)
		if result == false {
			return false
		}
	}

	return true
}
