// public.go contains the API routes that are public-facing.
package settings

import (
	"net/http"
    "github.com/labstack/echo/v4"
)

// PublicRoutes contains the public routes
func PublicRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {       
		return c.String(http.StatusOK, "Hello, World!\n")  
	})
}