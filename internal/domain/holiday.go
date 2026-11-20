// Package domain contains the core entities and value objects.
package domain

import "time"

// Holiday is a single entry of the Chilean holidays calendar.
type Holiday struct {
	// Date is the day of the holiday (UTC midnight).
	Date time.Time
	// Title is the holiday name.
	Title string
	// Type is the holiday category, e.g. "Civil" or "Religioso".
	Type string
	// Inalienable reports whether the holiday is legally inalienable (irrenunciable).
	Inalienable bool
	// Extra is free-form complementary information provided by the source.
	Extra string
}
