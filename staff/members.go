// members.go contains CRUD APIs for staff members.
package staff

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha-api/internal/methods"
)

// IndexStaff gets all staff members that fit the given parameters
// @ GET /admin/staff
func IndexStaff(c echo.Context) error {
	// Bind query parameters to model
	filter := new(methods.Staff)
	if err := c.Bind(filter); err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	// Get all applicable actions
	staffMembers, err := filter.Read(false)

	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// Return actions
	return c.JSON(http.StatusOK, staffMembers)
}

// CreateStaff creates a single Staff member
// @ POST /admin/staff/new
func CreateStaff(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	doAction := "staff-create"

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
	var staffDetails methods.Staff
	staffDetails.StaffRoleID, _ = strconv.Atoi(c.FormValue("role_id"))
	staffDetails.Username = c.FormValue("username")
	staffDetails.Email = c.FormValue("email")
	staffDetails.Password = c.FormValue("password")

	// Get all applicable actions
	err := staffDetails.Create()
	if err != nil {
		// Return the error
		staffLog.Create(false, methods.Error{Details: err,
			Data: staffDetails})
		c.Logger().Error(err)

		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// Return actions
	staffLog.Create(true, methods.Error{Data: staffDetails})
	return c.JSON(http.StatusOK, staffDetails)
}

// ReadStaff reads a single Staff member
// @ GET /admin/staff/:id
func ReadStaff(c echo.Context) error {
	// Bind query parameters to model
	var filter methods.Staff
	filter.ID, _ = strconv.Atoi(c.Param("id"))

	// Get all applicable actions
	staffMembers, err := filter.Read(true)

	if err != nil {
		// Return the error
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	staffMember := staffMembers[0]

	// Get Staff Role
	role := methods.StaffRole{ID: staffMember.StaffRoleID}
	roles, _ := role.Read()
	staffMember.StaffRole = roles[0]

	// Return actions
	return c.JSON(http.StatusOK, staffMember)
}

// UpdateStaff updates a single Staff member
// @ POST /admin/staff/:id
func UpdateStaff(c echo.Context) error {
	user := methods.CurrentAuthStaff(c.Get("user"))
	doAction := "staff-update"

	// Setup log
	staffLog := methods.StaffLog{
		StaffID: user.ID,
		StaffAction: methods.StaffAction{Name:doAction},
		IPAddress: c.RealIP(),
	}

	// Bind query parameters to model
	var staffMember methods.Staff
	staffMember.ID, _ = strconv.Atoi(c.Param("id"))

	// Check if user can edit other users (either with StaffPermission
	// or by editing themselves)
	if ! user.CanStaff(doAction) || user.ID == staffMember.ID {
		staffLog.Create(false, methods.Error{Message: "Unauthorised"})
		return echo.ErrUnauthorized
	}

	// Bind other parameters
	staffMember.Username = c.FormValue("username")
	staffMember.Email = c.FormValue("email")
	staffMember.Password = c.FormValue("password")
	staffMember.StaffRoleID, _ = strconv.Atoi(c.FormValue("role_id"))
	staffMember.IsActive, _ = strconv.ParseBool(c.FormValue("is_active"))

	// Get all applicable actions
	err := staffMember.Update(user.ID == staffMember.ID)
	if err != nil {
		// Return the error
		c.Logger().Error(err)
		staffLog.Create(false, 
			methods.Error{Details: err, Data: staffMember})

		return echo.NewHTTPError(http.StatusNotFound, staffMember)
	}

	// Return actions
	staffLog.Create(true, methods.Error{Data: staffMember})
	return c.JSON(http.StatusOK, staffMember)
}