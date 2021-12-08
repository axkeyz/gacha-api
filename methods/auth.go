// auth.go contains the structs (models) & methods to authenticate a user.
package methods

import (
	"time"
	"log"
	"os"
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"github.com/axkeyz/gacha/config"
)

// AUTH: An Auth struct stores the username & password used for
// authentication.
type Auth struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

// AUTH METHODS - Staff
// (auth *Auth).Staff authenticates a staff member using the Auth
// struct and staves the details in the JWTToken (token) pointer.
func (auth *Auth) Staff(token *JWTToken) (error) {
	// Setup database
	db := config.SetupDB()
	defer db.Close()

	q := `SELECT staff.id, staff.staff_role_id, staff_role.name, email FROM
	staff INNER JOIN staff_role on staff.staff_role_id = staff_role.id
	WHERE (lower(username) = lower($1) OR lower(email) = lower($1)) AND
	password = crypt($2, password)`

	// Check if user exists
	if err := db.QueryRow(q, auth.Username, auth.Password).Scan(
		&token.ID, &token.RoleID, &token.Role, &token.Email,
	); err != nil {
		return err
	} 

	return nil
}

// UpdateStaffLogin updates the last_login timestamp of a staff user on
// the staff table to the current timestamp.
func (auth *Auth) UpdateStaffLastLogin() {
	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// On authorised attempt, update last login
	if _, err := db.Exec(`Update staff set last_login = CURRENT_TIMESTAMP 
	WHERE username = $1 OR email = $1`, auth.Username); err != nil {
		log.Println(err)
	}
}

// TOKENS
// JWT Token struct stores the claims inside a JWT token.
type JWTToken struct {
	ID int `json:"id"`
	Name  string `json:"name"`
	Staff bool   `json:"staff"`
	RoleID int `json:"role_id"`
	Role string `json:"role"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// TOKENS - staff methods
// CreateStaffToken generates a new staff token, when given an Auth model.
// A staff token is only generated if the username / password combo matches
// the values stored in the database. This function returns the staff token
// as a string and an error if applicable.
func (token *JWTToken) CreateStaffToken(auth Auth) (string, error) {
	// Setup defaults for a staff JWT token
	token.Staff = true
	token.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()

	if err := auth.Staff(token); err != nil {
		// Authenticating the staff failed
		log.Println(err)
		return "", err
	}
	
	// Create token
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, token)

	// Generate encoded token and send it as response.
	encodedToken, err := t.SignedString([]byte(os.Getenv("APP_SECRET")));
	if err != nil {
		// Token encoding failed, return error
		return "", err
	} else {
		// Update staff login time
		auth.UpdateStaffLastLogin()
		// Return encoded token
		return encodedToken, nil
	}
}

// CurrentAuthStaff packs the currently authenticated staff member's
// JWT token into its readable form.
func CurrentAuthStaff(user interface{}) JWTToken {
	// Token data model
	token := JWTToken{}

	// Retrieve and JSON-ify data stored in token
	userData := user.(*jwt.Token)
	tmp, _ := json.Marshal(userData.Claims)
	_ = json.Unmarshal(tmp, &token)

	// Return a string of the user's role
	return token
}

func (user *JWTToken) Permissions() ([]StaffAction ) {
	var permissions []StaffAction
	var permission StaffAction

	// Setup database
	db := config.SetupDB()
	defer db.Close()

	// Setup query
	q := `SELECT staff_action_id, staff_action.name
	FROM staff_permission INNER JOIN staff_action ON 
	staff_permission.staff_action_id = staff_action.id where
	staff_role_id = $1`

	if rows, err := db.Query(q, user.RoleID); err == nil {
		
		// Map each row (of permission data) to pointers
		for rows.Next() {
			// Save data to pointer
			err = rows.Scan(&permission.ID, &permission.Name)

			if err != nil {
				// Display error if rows.Scan causes an error
				return permissions
			}

			// append data to permissions map
			// permissions[staff_action_name] = staff_action_id
			permissions = append(permissions, permission)
		}
	}
	// Return all permissions
	return permissions
}

// JWTToken.Can("action-name") returns true if the user can
// perform the identified action name.
func (user *JWTToken) Can(doAction string) ( bool ) {
	if ! user.Staff {
		return false
	}

	permissions := user.Permissions()

	for _, i := range permissions {
		if i.Name == doAction {
			return true
		}
	}
	return false
}