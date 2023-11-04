package domain

type TrackingType string

const (
	MoreThanValue  = "more than price"
	LowerThanValue = "lower than price"
	MoreThanMA     = "more than MA"
)

type Indicator struct {
	Id          int64
	AlertID     int64
	IndicatorID TrackingType
	Value       int64
}
