package domain

type TrackingType string

const (
	MoreThanValue  TrackingType = "more than price"
	LowerThanValue TrackingType = "lower than price"
	MoreThanMA     TrackingType = "more than MA"
)

type Indicator struct {
	Id          int64
	AlertID     int64
	IndicatorID TrackingType
	Value       int64
}
