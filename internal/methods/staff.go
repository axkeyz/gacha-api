// user.go contains structs that relate to staff (admin) users of
// gacha-api.
package methods

import (
	"strings"
	"strconv"
	"errors"

	"github.com/axkeyz/gacha-api/config"
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

//========================== STAFF ROLE =========================//
// StaffRole involves the staff roles of the API.
//
// Every staff role has a different name. Each staff user can
// only have one staff role. Staff roles affect a staff user's
// permissions.
//
// This is directly mapped to the staff_role table. 
//===============================================================//

type StaffRole struct {
	ID int `query:"id" json:",omitempty"`
	Name string `query:"name" json:",omitempty"`
	Staff []Staff `json:",omitempty"`
	StaffPermission []StaffPermission `json:",omitempty"`
	Pagination
}

// StaffRole.Create creates a new StaffRole in the database 
// given a unique name.
func (role *StaffRole) Create() error {
	if role.Name == "" {
		return errors.New("Role name cannot be empty")
	}

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	query := `INSERT INTO staff_role (name) VALUES ($1)`

	if _, err := db.Exec(query, role.Name); err != nil {
		return err
	}

	return nil
}

// StaffRole.Read returns a struct of all StaffRole(s) that
// fit the given filter *StaffRole.
func (filter *StaffRole) Read() ([]StaffRole, error) {
	var roles []StaffRole
	var role StaffRole

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	main := `SELECT id, name FROM staff_role`
	where := filter.Filter()
	sort := filter.Pagination.Query()

	if rows, err := db.Query(main+where+sort); err != nil {
		return roles, err
	} else {
		for rows.Next() {
			// Save data to pointer
			err = rows.Scan(&role.ID, &role.Name)

			if err != nil {
				// Display error if rows.Scan causes an error
				return roles, err
			}

			roles = append(roles, role)
		}
	}
	return roles, nil
}

// StaffRole.Update updates the name of a StaffRole given
// its ID.
func (role *StaffRole) Update() error {
	// Check that the name is not empty
	if role.Name == "" {
		return errors.New("Role name cannot be empty")
	}

	// Setup database & query
	db := config.SetupDB()
	defer db.Close()

	if _, err := db.Exec( 
		"UPDATE staff_role SET name = $1 WHERE id = $2", 
		role.Name, role.ID,
	); err == nil {
		// Return nothing
		return nil
	} else {
		// Return error
		return err
	}
}

func (role *StaffRole) Filter() string {
	var items []string
	
	if role.Name != "" {
		items = append(items, "lower(name) LIKE lower('%"+role.Name+"%')")
	}

	if role.ID != 0 {
		items = append(items, "id = "+strconv.Itoa(role.ID))
	}

	if len(items) > 0 {
		return " WHERE " + strings.Join(items, " AND ")
	} else {
		return ""
	}
}

func (role *StaffRole) GetStaff() ([]Staff, error) {
	var staffMembers []Staff
	var staffMember Staff

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	rows, err := db.Query(`SELECT id, username, email, is_active
	FROM staff WHERE staff_role_id = $1`, role.ID)

	if err != nil {
		return staffMembers, err
	} else {
		// Map each row (of permission data) to pointers
		for rows.Next() {
			// Save data to pointer
			err = rows.Scan(&staffMember.ID, &staffMember.Username,
			&staffMember.Email, &staffMember.IsActive)

			if err != nil {
				// Display error if rows.Scan causes an error
				return staffMembers, err
			}

			// append data to staffMembers
			staffMembers = append(staffMembers, staffMember)
		}
	}
	return staffMembers, nil
}

//========================= STAFF ACTION =========================//
// StaffAction are the generic staff actions, that are used to name
// shared staff role permissions.
//
// Every staff role has different permissions, which defines whether
// they can perform certain actions or not.
//
// This is directly mapped to the staff_action table. 
//================================================================//

type StaffAction struct {
	ID int `query:"id" json:",omitempty"`
	Name string `query:"name" json:",omitempty"`
	StaffRole []StaffRole
	StaffPermission []StaffPermission
	Pagination
}

// StaffAction.Index returns all staff actions that fit the
// criteria of *StaffAction.
func (filter *StaffAction) Index() ([]StaffAction, error) {
	var action StaffAction
	var actions []StaffAction

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	main := `SELECT id, name FROM staff_action`
	where := filter.Filter(true)
	sort := filter.Pagination.Query()

	if rows, err := db.Query(main+where+sort); err != nil {
		return actions, err
	} else {
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

// StaffAction.Create creates a single staff action given
// a name.
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

// StaffAction.Read returns the result of a single staff
// action given its ID.
func (filter *StaffAction) Read() (StaffAction, error) {
	var action StaffAction

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	main := `SELECT id, name FROM staff_action WHERE id = $1`
	sort := filter.Pagination.Query()

	err := db.QueryRow(main+sort, filter.ID).Scan(
		&action.ID, &action.Name,
	)

	// Get related permissions
	filterStaffPermission := StaffPermission{
		StaffActionID: action.ID,
	}

	action.StaffPermission, err = filterStaffPermission.Read()

	// Return row
	return action, err
}

// StaffAction.Update updates the name of a staff action 
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

// StaffAction.Filter generates an SQL "WHERE" clause, given
// the filter (*StaffAction). When using notExplicit = true,
// this function generates a WHERE clause with LIKE = '%val%'
// instead of direct =.
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

//======================= STAFF PERMISSION =======================//
// StaffPermission involves the staff permissions of the API.
//
// Every staff role has different permissions, which defines whether
// they can perform certain actions or not.
//
// This is directly mapped to the staff_action table. 
//===============================================================//

type StaffPermission struct {
	ID int `query:"id" json:",omitempty"`
	StaffRoleID int `query:"role_id" json:",omitempty"`
	StaffRole StaffRole `json:",omitempty"`
	StaffActionID int `query:"action_id" json:",omitempty"`
	StaffAction StaffAction `json:",omitempty"`
	Pagination
}

// StaffPermission.Create creates a StaffPermission in the database
// given that the permission *StaffPermission's StaffRoleID and
// StaffActionID is not 0.
func (permission *StaffPermission) Create() error {
	if permission.StaffRoleID == 0 || permission.StaffActionID == 0 {
		return errors.New("Role ID and Action ID cannot be empty")
	}

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	query := `INSERT INTO staff_permission (staff_role_id, 
		staff_action_id) VALUES ($1, $2)`

	if _, err := db.Exec(query, permission.StaffRoleID, 
		permission.StaffActionID); err != nil {
		return err
	}

	return nil
}

// StaffPermission.Read returns all StaffPermissions that fulfil
// all parameters of the filter *StaffPermission.
func (filter *StaffPermission) Read() ([]StaffPermission, error) {
	var permissions []StaffPermission
	var permission StaffPermission

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	main := `SELECT staff_permission.id as permission_id, staff_role.id as role_id, 
	staff_role.name as role_name, staff_action_id, staff_action.name as action_name 
	FROM staff_permission INNER JOIN staff_action ON staff_permission.staff_action_id
	= staff_action.id INNER JOIN staff_role ON staff_role.id = 
	staff_permission.staff_role_id`
	where := filter.Filter()
	sort := filter.Pagination.Query()

	if rows, err := db.Query(main+where+sort); err == nil {
		
		// Map each row (of permission data) to pointers
		for rows.Next() {
			// Save data to pointer
			err = rows.Scan( &permission.ID, 
				&permission.StaffRole.ID, &permission.StaffRole.Name,
				&permission.StaffAction.ID, &permission.StaffAction.Name,
			)

			if err != nil {
				// Display error if rows.Scan causes an error
				return permissions, err
			}

			// append data to permissions
			permissions = append(permissions, permission)
		}
	}

	// Return all permissions
	return permissions, nil
}

// StaffPermission.Update updates a StaffPermission in the database
// given the permission.ID.
func (permission *StaffPermission) Update() error {
	// Check that the ids are not empty
	if permission.ID == 0 || permission.StaffRoleID == 0 || 
		permission.StaffActionID == 0 {
		return errors.New("Role ID and Action ID cannot be empty")
	}

	// Setup database & query
	db := config.SetupDB()
	defer db.Close()

	if _, err := db.Exec( 
		`UPDATE staff_permission SET staff_role_id = $1, 
		staff_action_id = $2 WHERE id = $3`, 
		permission.StaffRoleID, permission.StaffActionID, permission.ID,
	); err == nil {
		// Return nothing
		return nil
	} else {
		// Return error
		return err
	}
}

// StaffPermission.Delete deletes a StaffPermission in the database
// given the permission.ID.
func (permission *StaffPermission) Delete() error {
	// Check that the ids are not empty
	if permission.ID == 0 {
		return errors.New("Permission ID cannot be empty")
	}

	// Setup database & query
	db := config.SetupDB()
	defer db.Close()

	if _, err := db.Exec( 
		`DELETE FROM staff_permission WHERE id = $1`, permission.ID,
	); err == nil {
		// Return nothing
		return nil
	} else {
		// Return error
		return err
	}
}

// StaffPermission.Filter is a helper function that generates a
// filter string ("WHERE") given the filter parameters.
func (filter *StaffPermission) Filter() string {
	var items []string

	// Add parameters that are not null
	if filter.ID != 0 {
		items = append(items, "staff_action.id = "+strconv.Itoa(filter.ID))
	}
	if filter.StaffRoleID != 0 {
		items = append(items, "staff_action.id = "+strconv.Itoa(filter.StaffRoleID))
	}
	if filter.StaffActionID != 0 {
		items = append(items, "staff_action.id = "+strconv.Itoa(filter.StaffActionID))
	}

	if len(items) > 0 {
		return " WHERE " + strings.Join(items, " AND ")
	} else {
		return ""
	}
}