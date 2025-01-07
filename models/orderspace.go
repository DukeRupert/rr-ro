// models/orderspace.go
package models

type OrderspaceResponse struct {
	HasMore  bool                `json:"has_more"`
	Products []OrderspaceProduct `json:"products"`
}

type OrderspaceProduct struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Code           string           `json:"code"`
	Description    string           `json:"description"`
	Images         []string         `json:"images"` // Array of image URLs
	Active         bool             `json:"active"`
	Categories     []Category       `json:"categories"`
	Variants       []ProductVariant `json:"product_variants"`
	VariantOptions []string         `json:"variant_options"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductVariant struct {
	ID              string            `json:"id"`
	SKU             string            `json:"sku"`
	Options         map[string]string `json:"options"`
	PriceListPrices []PriceList       `json:"price_list_prices"`
	UnitPrice       float64           `json:"unit_price"`
	Backorder       bool              `json:"backorder"`
}

type PriceList struct {
	ID        string  `json:"id"`
	UnitPrice float64 `json:"unit_price"`
}
