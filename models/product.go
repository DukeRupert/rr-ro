// models/product.go
package models

type Product struct {
    ID       string
    Name     string
    ImageURL string
    Variants []Variant
}

type Variant struct {
    SKU       string
    Size      string
    Grind     string
    Price     float64
    Available bool
}