package staff

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha/internal/methods"
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

// CreateStaffAction creates a new staff action/permission when
// given the name of the action @ POST /admin/actions/new
func CreateStaffAction(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	
	if user.CanStaff("staffaction-create") {
		action := methods.StaffAction{
			Name: c.FormValue("name"),
		}

		if err := action.Create(); err != nil {
			c.Logger().Error(err)
			return echo.ErrBadRequest
		}

		// Return action
		return c.JSON(http.StatusOK, action)
	} else {
		return echo.ErrUnauthorized
	}
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

// UpdateStaffAction updates the name of a single StaffAction
// when given the StaffAction's id (int)
// @ POST /admin/actions/:id
func UpdateStaffAction(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	
	if user.CanStaff("staffaction-update") {
		var action methods.StaffAction

		// ID param is an integer
		action.ID, _ = strconv.Atoi(c.Param("id"))
		action.Name = c.FormValue("name")

		// Get all applicable actions
		err := action.Update()
		if err != nil {
			// Return the error
			c.Logger().Error(err)
			return echo.ErrBadRequest
		}

		// Return actions
		return c.JSON(http.StatusOK, action)
	} else {
		return echo.ErrUnauthorized
	}
}
