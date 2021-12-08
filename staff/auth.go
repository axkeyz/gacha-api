// auth.go creates an API that allows authentication of staff
// members.
package staff

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/axkeyz/gacha/methods"
)

// AuthenticateStaff attempts to authenticate a user with the given
// username & password. @POST /admin/login
// If the authentication is a success, a JWT token is generated. If 
// authentication failed, ErrUnauthorized is returned.
func AuthenticateStaff(c echo.Context) error {
	// Create auth credentials and token
	auth := methods.Auth{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}
	token := methods.JWTToken{}

	// Check if user exists
	if encodedToken, err := token.CreateStaffToken(auth); err != nil {
		// Unauthorised login attempt
		c.Logger().Error(err)
		return echo.ErrUnauthorized
	} else {
		return c.JSON(http.StatusOK, echo.Map{
			"token": encodedToken,
		})
	}
}