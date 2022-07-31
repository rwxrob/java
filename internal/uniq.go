package internal

import (
	"fmt"
	"time"
)

// Isonan returns the GMT current time in ISO8601 (RFC3339) but for
// nanoseconds without any punctuation or the T.  This is frequently
// a very good unique suffix that has the added advantage of being
// chronologically sortable and more readable than the epoch and
// provides considerably more granularity than just Second.
func Isonan() string {
	t := time.Now()
	return fmt.Sprintf("%v%v",
		t.In(time.UTC).Format("20060102150405"),
		t.In(time.UTC).Format(".999999999")[1:],
	)
}
