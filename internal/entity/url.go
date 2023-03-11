package entity

import "time"

const UrlTableName = "urls"

// URL defines the domain model
type URL struct {
	Token     string
	TargetURL string
	Visits    int
	CreatedAt time.Time
}
