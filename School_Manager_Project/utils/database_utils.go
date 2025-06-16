package utils

import (
	"net/http"
	"strings"
)

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func AddSortFilters(r *http.Request, query string) (string, error) {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				return query, InvalidSortParameterError
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				return query, InvalidSortParameterError
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query, nil
}

func AddSearchFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, _ := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + param + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}
