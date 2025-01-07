// handlers/template.go
package handlers

import (
	"net/http"
	"orderspace-integration/models"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) HandleTemplateSubmission(c echo.Context) error {
    // Parse form values
    quantities := make(map[string]int)
    for key, values := range c.Request().Form {
        if strings.HasPrefix(key, "quantities[") && len(values) > 0 {
            sku := strings.TrimSuffix(strings.TrimPrefix(key, "quantities["), "]")
            qty, err := strconv.Atoi(values[0])
            if err != nil {
                continue
            }
            if qty > 0 {
                quantities[sku] = qty
            }
        }
    }

    // Parse dates and frequency
    initialDate, err := time.Parse("2006-01-02", c.FormValue("initial_date"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid date format",
        })
    }

    freqNum, err := strconv.Atoi(c.FormValue("frequency_number"))
    if err != nil || freqNum < 1 {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid frequency number",
        })
    }

    freqUnit := c.FormValue("frequency_unit")
    if freqUnit != "weeks" && freqUnit != "months" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid frequency unit",
        })
    }

    // Create template
    template := models.OrderTemplate{
        InitialDate:   initialDate,
        FrequencyNum:  freqNum,
        FrequencyUnit: freqUnit,
    }

    // Add items to template
    for sku, qty := range quantities {
        template.Items = append(template.Items, models.TemplateItem{
            SKU:      sku,
            Quantity: qty,
        })
    }

    // TODO: Save template to database or send to Orderspace API

    return c.JSON(http.StatusOK, template)
}