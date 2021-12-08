// user.go contains structs that relate to staff (admin) users of
// gacha-api.
package methods

import (
	"strings"
	"strconv"
	"errors"

	"github.com/axkeyz/gacha/config"
)

// A staff struct stores the data of a staff members.
type Staff struct {
	ID int `json:",omitempty"`
	StaffRoleID int `json:",omitempty"`
	StaffRole StaffRole `json:",omitempty"`
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
	Email string `json:",omitempty"`
	IsActive bool `json:",omitempty"`
	JoinedAt string `json:",omitempty"`
	LastLogin string `json:",omitempty"`
}

type StaffRole struct {
	ID int `json:",omitempty"`
	Name string `json:",omitempty"`
	Staff []Staff `json:",omitempty"`
	StaffAction []StaffAction `json:",omitempty"`
}

type StaffAction struct {
	ID int `query:"id" json:",omitempty"`
	Name string `query:"name" json:",omitempty"`
	Pagination
}

func (filter *StaffAction) Index() ([]StaffAction, error) {
	var action StaffAction
	var actions []StaffAction

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	main := `SELECT id, name FROM staff_action`
	where := filter.Filter(true)
	sort := filter.Pagination.Query("staff_action")

	if rows, err := db.Query(main+where+sort); err == nil {
		
		// Map each row (of permission data) to pointers
		for rows.Next() {
			// Save data to pointer
			err = rows.Scan(&action.ID, &action.Name)

			if err != nil {
				// Display error if rows.Scan causes an error
				return actions, err
			}

			// append data to permissions map
			// permissions[staff_action_name] = staff_action_id
			actions = append(actions, action)
		}
	}

	// Return all permissions
	return actions, nil
}

func (action *StaffAction) Create() (error) {
	if action.Name == "" {
		return errors.New("action.Name cannot be empty")
	}

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	query := `INSERT INTO staff_action (name) VALUES ($1)`

	if _, err := db.Exec(query, action.Name); err != nil {
		return err
	}

	return nil
}

func (filter *StaffAction) Read() (StaffAction, error) {
	var action StaffAction

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	main := `SELECT id, name FROM staff_action`
	where := filter.Filter(false)
	sort := filter.Pagination.Query("staff_action")

	err := db.QueryRow(main+where+sort).Scan(
		&action.ID, &action.Name,
	)

	// Return row
	return action, err
}

// Update StaffAction updates the name of a staff action 
// given its ID.
func (action *StaffAction) Update() (error) {
	// Check that the name is not empty
	if action.Name == "" {
		return errors.New("action.Name cannot be empty")
	}

	// Setup database & query
	db := config.SetupDB()
	defer db.Close()

	if _, err := db.Exec( 
		"UPDATE staff_action SET name = $1 WHERE id = $2", 
		action.Name, action.ID,
	); err == nil {
		// Return nothing
		return nil
	} else {
		// Return error
		return err
	}
}

func (filter *StaffAction) Filter(notExplicit bool) (string) {
	var items []string
	f1 := " = '"
	f2 := "'"

	if notExplicit {
		f1 = "::text LIKE '%"
		f2 = "%'"
	}
	
	if filter.Name != "" {
		items = append(items, 
			"staff_action.name"+f1+filter.Name+f2)
	}

	if filter.ID != 0 {
		items = append(items, 
			"staff_action.id"+f1+strconv.Itoa(filter.ID)+f2)
	}

	if len(items) > 0 {
		return " WHERE " + strings.Join(items, " AND ")
	} else {
		return ""
	}
}