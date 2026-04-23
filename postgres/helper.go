package postgres

import (
	"fmt"
	"strings"
	"time"
)

func toAge(createdAt time.Time) string {
	age := int(time.Since(createdAt).Hours() / 24)
	return fmt.Sprintf("%dd", age)
}

func redactValue(value string) string {
	// Redact 90% of its length
	percent := 90.0 / 100.0
	stopRedacted := len(value) - int(float64(len(value))*percent)

	return fmt.Sprintf("%sxxxxxxxxx", value[:stopRedacted])
}

func ensureSort(sort *string) string {
	if sort == nil {
		return "ASC"
	}

	strSort := strings.ToUpper(*sort)
	if strSort != "DESC" && strSort != "ASC" {
		strSort = "ASC"
	}
	return strSort
}
