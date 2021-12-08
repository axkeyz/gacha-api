package staff

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha/internal/methods"
)

func IndexStaffPermissions(c echo.Context) error {
	// Bind query parameters to model
	filterStaffPermissions := new(methods.StaffPermission)
	if err := c.Bind(filterStaffPermissions); err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Get all applicable actions
	permissions, err := filterStaffPermissions.Read()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, permissions)
}

func CreateStaffPermission(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	if ! user.CanStaff("staffpermission-create") {
		return echo.ErrUnauthorized
	}

	// Bind query parameters to model
	action_id, _ := strconv.Atoi(c.FormValue("action_id"))
	role_id, _ := strconv.Atoi(c.FormValue("role_id"))

	permissions := methods.StaffPermission{
		StaffActionID: action_id,
		StaffRoleID: role_id,
	}

	// Get all applicable actions
	err := permissions.Create()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, permissions)
}

func ReadStaffPermission(c echo.Context) error {
	// Bind query parameters to model
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	} 
	
	filterStaffPermissions := methods.StaffPermission{
		ID: id,
	}

	// Get all applicable actions
	permissions, err := filterStaffPermissions.Read()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, permissions[0])
}

func UpdateStaffPermission(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	if ! user.CanStaff("staffpermission-update") {
		return echo.ErrUnauthorized
	}

	// Bind query parameters to model
	id, _ := strconv.Atoi(c.Param("id"))
	action_id, _ := strconv.Atoi(c.FormValue("action_id"))
	role_id, _ := strconv.Atoi(c.FormValue("role_id"))

	permissions := methods.StaffPermission{
		ID: id,
		StaffActionID: action_id,
		StaffRoleID: role_id,
	}

	// Get all applicable actions
	err := permissions.Update()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, permissions)
}

func DeleteStaffPermission(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	if ! user.CanStaff("staffpermission-delete") {
		return echo.ErrUnauthorized
	}

	// Bind query parameters to model
	id, _ := strconv.Atoi(c.Param("id"))

	permissions := methods.StaffPermission{
		ID: id,
	}

	// Get all applicable actions
	err := permissions.Delete()
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Return actions
	return c.JSON(http.StatusOK, permissions)
}