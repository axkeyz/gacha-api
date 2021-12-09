package staff

import (
	"net/http"
	
	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha-api/internal/methods"
)

func TestAuthenticateStaff(c echo.Context) error {
	// Get the currently authenticated staff member
	user := methods.CurrentAuthStaff(c.Get("user"))

	return c.JSON(http.StatusOK, user)
}

func TestStaffPermission(c echo.Context) error {
	// Get the currently authenticated staff member
	user := methods.CurrentAuthStaff(c.Get("user"))
	
	permission := user.CanStaff("staff-update")

	return c.JSON(http.StatusOK, permission)
}