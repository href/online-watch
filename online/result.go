package online

import (
	"time"
)

// Result contains the details of a check
type Result struct {
	// Time denotes the time the check completed
	Done time.Time

	// Took contains how long it took
	Took time.Duration

	// Okay is true if the check succeeded
	Okay bool

	// Consecutive is the number of results that came before with the same
	// result, including this one (i.e., it starts at 1).
	Consecutive int

	// Duration since the last time the result changed to ok, or from ok.
	Duration time.Duration

	// Hints contains a textual representation of the result
	Hints []string

	// Target identifies the target this result was run on
	Target *Target
}

// Text shows the result as single character
func (o Result) ShortText() string {
	if o.Okay {
		return "✔︎"
	} else {
		return "✖︎"
	}
}
