package go_template

import (
	"errors"
	"time"
)

type PlainTime time.Time

func (tm PlainTime) MarshalJSON() ([]byte, error) {
	t := time.Time(PlainTime{})
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.RFC3339Nano)
	b = append(b, '"')
	return b, nil
}
