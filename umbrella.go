package umbrella

import "time"

var (
	processStartTime time.Time
)

func init() {
	processStartTime = time.Now()
}
