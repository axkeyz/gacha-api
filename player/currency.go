// currency.go tracks player currency transactions. It also allows
// issuing of player currency.
package player

import (
	"net/http"

	"github.com/axkeyz/gacha-api/internal/methods"
	"github.com/labstack/echo/v4"
)

func IndexPlayerCurrency(c echo.Context) error {
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
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	// Return actions
	return c.JSON(http.StatusOK, actions)
}
