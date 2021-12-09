package methods

import (
	"fmt"
	"strings"
)

type Pagination struct {
	Offset int  `query:"offset" json:",omitempty"`
    Limit  int  `query:"limit" json:",omitempty"`
    SortBy string `query:"sortby" json:",omitempty"`
    SortDir string `query:"dir" json:",omitempty"`
}

func (filters *Pagination) Query() string {
	order := ""

	if filters.SortBy != "" {
		// Default direction to ascending
		direction := "asc"
		if filters.SortDir == "desc" {
			// Change directon to descending
			direction = "desc"
		}
		order = fmt.Sprintf(" ORDER BY %s %s ", 
			filters.SortBy, direction)

		// Pagination only possible if there is a sortby & direction
		if filters.Limit != 0 {
			// Default offset to 0
			offset := 0
			if filters.Offset != 0 {
				offset = filters.Offset
			}
			// Update order clause to reflect pagination
			order = order + fmt.Sprintf("LIMIT %d OFFSET %d", filters.Limit, offset)
		}
	}

	return order
}

func GenerateWhereQuery(s map[string]interface{}, explicit bool) string {
	var params []string
	var param, f1, f2 string
	// var test bool

	if explicit {
		f1 = " = '"
		f2 = "'"
	} else {
		f1 = " LIKE lower('%%"
		f2 = "%%')"
	}

	for key, val := range s {
		param = fmt.Sprintf("%v", val)

		if IsNotPagination(key) && ! strings.Contains(param, "map") {
			params = append(params,
				fmt.Sprintf("lower(%s::text) %s%s%s", key, f1, param, f2))
		}
	}

	if len(params) > 0 {
		return " WHERE " + strings.Join(params, " AND ")
	} else {
		return ""
	}
}

// IsNotPagination returns true if the given key is not a pagination
// key.
func IsNotPagination(key string) bool {
	pagination := []string{"SortBy", "Asc", "Limit", "Offset"}

	for _, p := range pagination {
		if strings.Contains(key, p) {
			return false
		}
	}

	return true
}

type Error struct {
	Code string `json:",omitempty"`
	Message string `json:",omitempty"`
	Details interface{} `json:",omitempty"`
	Data interface{} `json:",omitempty"`
}