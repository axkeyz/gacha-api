package main

import (
	"os"
	"net/http"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/joho/godotenv"

	"github.com/axkeyz/gacha/methods"
	"github.com/axkeyz/gacha/staff"
)

func main() {
	// Open .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	// API groups - restricted for staff only
	a := e.Group("/admin")
	t := e.Group("/test")

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	// JWT-middleware
	adminJWT := middleware.JWTConfig{
		Claims:     &methods.JWTToken{},
		SigningKey: []byte(os.Getenv("APP_SECRET")),
	}
	a.Use(middleware.JWTWithConfig(adminJWT))
	t.Use(middleware.JWTWithConfig(adminJWT))

	// Set IP address extractor
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	// Routes
	e.GET("/", func(c echo.Context) error {       
		return c.String(http.StatusOK, "Hello, World!\n")  
	})

	e.POST("/admin/login", staff.AuthenticateStaff)

	// Tests
	t.GET("/admin/login", staff.TestAuthenticateStaff)
	t.GET("/admin/actions/can", staff.TestStaffPermission)
	
	a.GET("/actions", staff.IndexStaffActions)
	a.POST("/actions/new", staff.CreateStaffAction)
	a.GET("/actions/:id", staff.ReadStaffAction)
	a.POST("/actions/:id", staff.UpdateStaffAction)

	e.Logger.Fatal(e.Start(":1588"))
}