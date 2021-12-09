package main

import (
	"os"
	"net/http"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/joho/godotenv"

	"github.com/axkeyz/gacha-api/internal/methods"
	// "github.com/axkeyz/gacha-api/internal/utils"
	"github.com/axkeyz/gacha-api/staff"
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

	// loggerConfig := middleware.DefaultLoggerConfig

	// f, err := os.OpenFile("logs/api-"+utils.TodayIs()+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal("Could not open log file, (%s)", err.Error())
	// 	os.Exit(-1)
	// }

	// defer f.Close()

	// loggerConfig.Output = f

	// // Middlewares
	// e.Use(middleware.LoggerWithConfig(loggerConfig))
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

	a.GET("/permissions", staff.IndexStaffPermissions)
	a.POST("/permissions/new", staff.CreateStaffPermission)
	a.GET("/permissions/:id", staff.ReadStaffPermission)
	a.POST("/permissions/:id", staff.UpdateStaffPermission)
	a.DELETE("/permissions/:id", staff.DeleteStaffPermission)

	a.GET("/roles", staff.IndexStaffRoles)
	a.POST("/roles/new", staff.CreateStaffRole)
	a.GET("/roles/:id", staff.ReadStaffRole)
	a.POST("/roles/:id", staff.UpdateStaffRole)

	e.Logger.Fatal(e.Start(":1588"))
}