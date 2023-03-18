package entity

import "time"

const UrlTableName = "url"

// Url defines the domain model
type Url struct {
	Token     string    `db:"token"`
	TargetUrl string    `db:"target_url"`
	Visits    int       `db:"visits"`
	CreatedAt time.Time `db:"created_at"`
}
