package utils

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