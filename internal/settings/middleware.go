// middleware.go contains the middlewares of the API.
package settings

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/axkeyz/gacha-api/internal/methods"
)

// RegisterMiddlewares includes most of the important middlewares
func RegisterMiddlewares(e *echo.Echo) {
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
}

// adminJWT & AdminMiddleware create the JWT-Token-based middleware
// used to authenticate staff on admin/staff routes.
var adminJWT = middleware.JWTConfig{
	Claims:     &methods.JWTToken{},
	SigningKey: []byte(os.Getenv("APP_SECRET")),
}

func AdminMiddleware(g *echo.Group) {
	g.Use(middleware.JWTWithConfig(adminJWT))
}