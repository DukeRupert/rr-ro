// models/template.go
package models

import "time"

type OrderTemplate struct {
    Items          []TemplateItem
    InitialDate    time.Time
    FrequencyNum   int
    FrequencyUnit  string // "weeks" or "months"
}

type TemplateItem struct {
    SKU      string
    Quantity int
    Price    float64
}