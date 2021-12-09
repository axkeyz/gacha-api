// roles.go contains functions that involve the CRUD of
// staff roles.
package staff

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha-api/internal/methods"
)

// IndexStaffActions returns a list of all staff roles available
// @ GET /admin/roles
func IndexStaffRoles(c echo.Context) error {
	// Bind query parameters to model
	filter := new(methods.StaffRole)
	if err := c.Bind(filter); err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Get all applicable actions
	roles, err := filter.Read()

	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// Return actions
	return c.JSON(http.StatusOK, roles)
}

// CreateStaffRoles creates a new staff role
// @ POST /admin/roles/new
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

// ReadStaffRole reads a single staff role
// @ GET /admin/roles/:id
func ReadStaffRole(c echo.Context) error {
	// Bind query parameters to model
	var filter methods.StaffRole
	filter.ID, _ = strconv.Atoi(c.Param("id"))

	// Get all applicable actions
	roles, err := filter.Read()
	role := roles[0]
	staff, _ := role.GetStaff()

	filterStaffPermissions := methods.StaffPermission{
		StaffRoleID: filter.ID,
	}

	role.StaffPermission, _ = filterStaffPermissions.Read()


	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	role.Staff = staff

	// Return actions
	return c.JSON(http.StatusOK, role)
}

// UpdateStaffRole reads a single staff role
// @ POST /admin/roles/:id
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