package staff

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha/methods"
)

// IndexStaffActions returns a list of all staff actions available
// @ GET /admin/actions
func IndexStaffActions(c echo.Context) error {
	// Bind query parameters to model
	filterStaffActions := new(methods.StaffAction)
	if err := c.Bind(filterStaffActions); err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Get all applicable actions
	actions, err := filterStaffActions.Index()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, actions)
}

// ReadStaffAction returns the details of a single StaffAction
// when given the StaffAction's id (int) or exact name (string)
// @ GET /admin/actions/:id
func ReadStaffAction(c echo.Context) error {
	var filterStaffActions methods.StaffAction

	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		// ID param is an integer
		filterStaffActions.ID = id
	} else {
		// ID param is a string, possible the name
		filterStaffActions.Name = c.Param("id")
	}

	// Get all applicable actions
	actions, err := filterStaffActions.Read()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, actions)
}

