// main.go
package main

import (
	"os"

	"orderspace-integration/handlers"
	"orderspace-integration/services"
	"orderspace-integration/template"

	"github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Configure zerolog
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	// Initialize Echo instance
	e := echo.New()

	// Initialize services
	orderSpaceService := services.NewOrderSpaceService(
		os.Getenv("ORDERSPACE_CLIENT_ID"),
		os.Getenv("ORDERSPACE_CLIENT_SECRET"),
	)

	// Initialize template renderer
	tmpl := template.NewTemplate()

	// Initialize handlers
	handlers := handlers.NewHandlers(orderSpaceService, tmpl.Templates)

	// Set up renderer
	e.Renderer = tmpl

	// Routes
	e.GET("/", handlers.HandleProducts)
	e.POST("/submit-template", handlers.HandleTemplateSubmission)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
