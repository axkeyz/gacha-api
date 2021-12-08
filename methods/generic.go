package methods

import (
	"fmt"
)

type Pagination struct {
	Offset int  `query:"offset" json:",omitempty"`
    Limit  int  `query:"limit" json:",omitempty"`
    SortBy string `query:"sortby" json:",omitempty"`
    SortDir string `query:"dir" json:",omitempty"`
}

func (filters *Pagination) Query(table string) string {
	order := ""

	if filters.SortBy != "" {
		// Default direction to ascending
		direction := "asc"
		if filters.SortDir == "desc" {
			// Change directon to descending
			direction = "desc"
		}
		order = fmt.Sprintf(" ORDER BY %s.%s %s ", table, 
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