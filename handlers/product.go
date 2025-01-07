// handlers/products.go
package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"orderspace-integration/models"
)

type ProductData struct {
	Title    string
	Products []models.Product
}

func (h *Handlers) HandleProducts(c echo.Context) error {
	logger := log.With().
		Str("handler", "HandleProducts").
		Str("method", "GET").
		Logger()

	logger.Info().Msg("fetching products from Orderspace API")

	osProducts, err := h.orderSpaceService.GetProducts()
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch products from Orderspace API")
		return fmt.Errorf("failed to fetch products: %w", err)
	}

	logger.Info().
		Int("product_count", len(osProducts)).
		Msg("successfully fetched products")

	var products []models.Product
	for _, osProduct := range osProducts {
		product := models.Product{
			ID:       osProduct.ID,
			Name:     osProduct.Name,
			ImageURL: "", // We'll set this below
		}

		// Get first image if available
		if len(osProduct.Images) > 0 {
			product.ImageURL = osProduct.Images[0]
		}

		// Transform variants
		for _, osVariant := range osProduct.Variants {
			var price float64
			if len(osVariant.PriceListPrices) > 0 {
				price = osVariant.PriceListPrices[0].UnitPrice
			}

			variant := models.Variant{
				SKU:       osVariant.SKU,
				Size:      osVariant.Options["Size"],
				Grind:     osVariant.Options["Grind"],
				Price:     price,
				Available: true, // All products seem to be "backorder: true" so we'll mark as available
			}
			product.Variants = append(product.Variants, variant)
		}

		products = append(products, product)
	}

	logger.Info().
		Int("transformed_count", len(products)).
		Msg("finished transforming products")

	data := ProductData{
		Title:    "Products",
		Products: products,
	}

	logger.Info().Msg("attempting to render template")
	err = c.Render(http.StatusOK, "layout", data)
	if err != nil {
		logger.Error().Err(err).Msg("template rendering failed")
		return err
	}

	logger.Info().Msg("successfully rendered template")
	return nil
}
