package entity

import "time"

const UrlTableName = "urls"

// Url defines the domain model
type Url struct {
	Token     string
	TargetUrl string
	Visits    int
	CreatedAt time.Time
}
