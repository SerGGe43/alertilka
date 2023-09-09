package domain

type TrackingType string

const (
	MoreThanValue  TrackingType = "more than price"
	LowerThanValue TrackingType = "lower than price"
	MoreThanMA     TrackingType = "more than MA"
)

type Indicator struct {
	id          int
	alertID     int
	indicatorID TrackingType
	value       int
}
