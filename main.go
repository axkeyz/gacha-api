package main

import (
	"github.com/labstack/echo/v4"
	"github.com/joho/godotenv"

	"github.com/axkeyz/gacha-api/internal/settings"
)

func main() {
	// Open .env file
	_ = godotenv.Load()

	// Register echo (public) & groups (admin & testing)
	e := echo.New()
	a := e.Group("/admin")
	t := e.Group("/test")

	// register middleware
	settings.RegisterMiddlewares(e)
	settings.AdminMiddleware(a)
	settings.AdminMiddleware(t)

	// Set IP address extractor
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// Routes
	settings.PublicRoutes(e)

	settings.AuthenticateAdmin(e)
	settings.TestAdminRoutes(t)
	settings.AdminRoutes(a)

	e.Logger.Fatal(e.Start(":1588"))
}