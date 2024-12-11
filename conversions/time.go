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
