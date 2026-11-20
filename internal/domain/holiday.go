package domain

import "time"

type Holiday struct {
	Date time.Time

	Title string

	Type string

	Inalienable bool

	Extra string
}
