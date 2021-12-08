package staff

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha/internal/methods"
)

func CreateStaffRole(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	doAction := "staffrole-update"

	// Setup log
	staffLog := methods.StaffLog{
		StaffID: user.ID,
		StaffAction: methods.StaffAction{Name:doAction},
		IPAddress: c.RealIP(),
	}
	
	if ! user.CanStaff(doAction) {
		staffLog.Create(false, methods.Error{Message: "Unauthorised"})
		return echo.ErrUnauthorized
	}

	// Bind query parameters to model
	role := methods.StaffRole{
		Name: c.FormValue("name"),
	}

	// Get all applicable actions
	err := role.Create()
	if err != nil {
		// Return the error
		staffLog.Create(false, methods.Error{Details: err,
			Data: role})
		c.Logger().Error(err)

		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// Return actions
	staffLog.Create(true, methods.Error{Data: role})
	return c.JSON(http.StatusOK, role)
}

func UpdateStaffRole(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	doAction := "staffpermission-update"

	// Setup log
	staffLog := methods.StaffLog{
		StaffID: user.ID,
		StaffAction: methods.StaffAction{Name:doAction},
		IPAddress: c.RealIP(),
	}

	if ! user.CanStaff(doAction) {
		staffLog.Create(false, methods.Error{Message: "Unauthorised"})
		return echo.ErrUnauthorized
	}

	// Bind query parameters to model
	role := methods.StaffRole{
		Name: c.FormValue("name"),
	}

	role.ID, _ = strconv.Atoi(c.Param("id"))

	// Get all applicable actions
	err := role.Update()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		staffLog.Create(false, 
			methods.Error{Details: err, Data: role})

		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// Return actions
	staffLog.Create(true, methods.Error{Data: role})
	return c.JSON(http.StatusOK, role)
}