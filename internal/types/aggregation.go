package types

// AggResolution represents the resolution of aggregation.
type AggResolution string

const (
	// AggResolutionDay represents the resolution of aggregation by day.
	AggResolutionDay AggResolution = "DAY"

	// AggResolutionHour represents the resolution of aggregation by hour.
	AggResolutionHour AggResolution = "HOUR"

	// AggResolutionMinute represents the resolution of aggregation by minute.
	AggResolutionMinute AggResolution = "MINUTE"
)

// ToDuration converts the aggregation resolution to duration in seconds.
func (ar AggResolution) ToDuration() uint {
	switch ar {
	case AggResolutionDay:
		return 60 * 60 * 24
	case AggResolutionHour:
		return 60 * 60
	case AggResolutionMinute:
		return 60
	}
	return 0
}

// AggSubject represents the subject of aggregation.
type AggSubject string

const (
	// AggSubjectTxsCount represents the type of aggregation by transaction count.
	AggSubjectTxsCount AggSubject = "TXS_COUNT"

	// AggSubjectGasUsed represents the type of aggregation by gas used.
	AggSubjectGasUsed AggSubject = "GAS_USED"
)
