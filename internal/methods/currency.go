// currency.go contains CRUD functions for player currency
package methods

import (
	"errors"

	"github.com/axkeyz/gacha-api/config"
	"github.com/axkeyz/gacha-api/internal/utils"
)

//============================= Currency ===========================//
// Currency is the currency (resources) API.
//
// This is directly related to all in-game currency & resources
//
// This is directly mapped to the currency table.
//===============================================================//
type Currency struct {
	ID          int    `json:",omitempty"`
	Name        string `json:",omitempty"`
	URL         string `json:",omitempty"`
	Description string `json:",omitempty"`
	IsActive    bool   `json:",omitempty"`
	CreatedAt   string `json:",omitempty"`
	UpdatedAt   string `json:",omitempty"`
}

// Currency.Create creates a new Currency in the database given
// parameters in the *Currency struct.
func (currency *Currency) Create() {
	if utils.HasNoEmptyParams(
		[]string{currency.Name, currency.URL, currency.Description},
	) {
		// Setup database
		db := config.SetupDB()
		defer db.Close()

		// Setup query
		q := `INSERT into currency (name, url, description, is_active)
		VALUES ($1, $2, $3, $4);`

		db.Exec(q, currency.Name, currency.URL,
			currency.Description, currency.IsActive,
		)
	}
}

// Currency.Read reads all Currency in the database given
// parameters in the *Currency struct.
func (currency *Currency) Read() {
}

// Currency.Update updates a Currency in the database given
// parameters in the *Currency struct.
func (currency *Currency) Update() error {
	// Check that the name is not empty
	if utils.HasNoEmptyParams(
		[]string{currency.Name, currency.URL, currency.Description},
	) {
		// Setup database & query
		db := config.SetupDB()
		defer db.Close()

		if _, err := db.Exec(
			`UPDATE currency SET name = $1, url = $2, 
			description = $3, is_active = $4 WHERE id = $2`,
			currency.Name, currency.URL, currency.Description,
			currency.IsActive, currency.ID,
		); err == nil {
			// Return nothing
			return nil
		} else {
			// Return error
			return err
		}
	} else {
		return errors.New("Missing required parameters")
	}
}

// Currency.Deactivate deactivates a Currency in the database given
// parameters in the *Currency struct.
func (currency *Currency) Deactivate() error {
	// Check that the ids are not empty
	if currency.ID == 0 {
		return errors.New("Currency ID cannot be empty")
	}

	// Setup database & query
	db := config.SetupDB()
	defer db.Close()

	if _, err := db.Exec(
		`UPDATE currency SET is_active = $1 WHERE id = $2`,
		false, currency.ID,
	); err == nil {
		// Return nothing
		return nil
	} else {
		// Return error
		return err
	}
}

//============================= Currency ===========================//
// Player Currency is the currency owned by the player.
//
// This is directly related to a player's in-game currency & resources.
//
// This is directly mapped to the currency table.
//===============================================================//
type PlayerCurrency struct {
	PlayerID   int      `json:",omitempty"`
	Player     Player   `json:",omitempty"`
	CurrencyID int      `json:",omitempty"`
	Currency   Currency `json:",omitempty"`
	Amount     int      `json:",omitempty"`
}

// PlayerCurrency.Read reads all PlayerCurrency in the database given
// parameters in the *Currency struct.
func (currency *PlayerCurrency) Read() {
}

// PlayerCurrency.Update updates a PlayerCurrency in the database given
// parameters in the *Currency struct.
func (currency *PlayerCurrency) Update() {
}

// PlayerCurrency.Delete deactivates a PlayerCurrency in the database given
// parameters in the *Currency struct.
func (currency *PlayerCurrency) Delete() {
}

type PlayerTransaction struct {
	ID         int      `json:",omitempty"`
	PlayerID   int      `json:",omitempty"`
	Player     Player   `json:",omitempty"`
	CurrencyID int      `json:",omitempty"`
	Currency   Currency `json:",omitempty"`
	Change     int      `json:",omitempty"`
	CreatedAt  string   `json:",omitempty"`
}
