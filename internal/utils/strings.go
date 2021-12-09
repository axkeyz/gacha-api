package utils

import(
	"regexp"
	"encoding/json"
)

// HasNoEmptyParams returns true if none of the strings in the
// params slice is empty.
func HasNoEmptyParams(params []string) bool {
	for _,val := range params {
		if val == "" {
			return false
		}
	}
	return true
}

// IsEmail returns true if the given string appears to be an email
func IsEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}

// MapifyStruct converts a struct to a map[string]interface
func MapifyStruct(s interface{}) map[string]interface{} {   
    var inInterface map[string]interface{}

    inrec, _ := json.Marshal(s)
    json.Unmarshal(inrec, &inInterface)

    return inInterface
}