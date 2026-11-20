// Package domain contains the core entities and value objects.
package domain

import "time"

type Holiday struct {
	Date        time.Time
	Title       string
	Type        string
	Inalienable bool
	Description string
}
