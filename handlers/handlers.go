// handlers/handlers.go
package handlers

import (
	"html/template"
	"orderspace-integration/services"
)

type Handlers struct {
    orderSpaceService *services.OrderSpaceService
    templates        *template.Template
}

func NewHandlers(orderSpaceService *services.OrderSpaceService, templates *template.Template) *Handlers {
    return &Handlers{
        orderSpaceService: orderSpaceService,
        templates:        templates,
    }
}