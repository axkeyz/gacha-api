// staff.go contains the routes for staff/admins
package settings

import (
    "github.com/labstack/echo/v4"
	"github.com/axkeyz/gacha-api/staff"
)

// Public-facing routes for staff users
func AuthenticateAdmin(e *echo.Echo) {
	e.POST("/admin/login", staff.AuthenticateStaff)
}

// Non-public facing routes for staff users
func AdminRoutes(a *echo.Group)  {
    a.GET("/staff", staff.IndexStaff)
	a.POST("/staff/new", staff.CreateStaff)
	a.GET("/staff/:id", staff.ReadStaff)
	a.POST("/staff/:id", staff.UpdateStaff)
	
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
}

// Non public facing routes for staff, usually used for API
// testing purposes
func TestAdminRoutes(t *echo.Group)  {
	t.GET("/admin/login", staff.TestAuthenticateStaff)
	t.GET("/admin/actions/can", staff.TestStaffPermission)
}