package conversions

import (
	"time"

	"github.com/jackc/pgtype"
)

func TimeToTimestamp(t time.Time) (pgtype.Timestamp, error) {
	timestamp := &pgtype.Timestamp{}

	return *timestamp, timestamp.Set(t)
}

func MustTimeToTimestamp(t time.Time) pgtype.Timestamp {
	out, err := TimeToTimestamp(t)
	if err != nil {
		panic(err)
	}
	return out
}

func UnixToTimestamp(unix int64) (pgtype.Timestamp, error) {
	return TimeToTimestamp(time.Unix(unix, 0))
}

func MustUnixToTimestamp(unix int64) pgtype.Timestamp {
	out, err := UnixToTimestamp(unix)
	if err != nil {
		panic(err)
	}
	return out
}
