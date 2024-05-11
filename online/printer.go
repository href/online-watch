package online

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/go-color-term/go-color-term/coloring"
)

// Printer provides snazzy result printing
type Printer struct {
	Output              io.Writer
	Verbose             bool
	ConsecutiveFailures int
	Streak              map[string]*Streak
}

// Streak tracks how long something was online or offline
type Streak struct {
	LongestOutage        time.Duration
	LongestFailureStreak int
}

// Print outputs the current state
func (p *Printer) Print(r Result) {

	if p.Streak == nil {
		p.Streak = make(map[string]*Streak)
	}

	title := r.Target.Title()

	if _, ok := p.Streak[title]; !ok {
		p.Streak[title] = &Streak{}
	}

	if !r.Okay {
		streak := p.Streak[title]
		streak.LongestFailureStreak = max(
			streak.LongestFailureStreak, r.Consecutive)
		streak.LongestOutage = max(
			streak.LongestOutage, r.Duration)
	}

	if !p.Verbose && (r.Okay || r.Consecutive < p.ConsecutiveFailures) {
		return
	}

	fmt.Fprint(p.Output, r.Done.Format("2006-01-02 15:04:05.000 "))
	fmt.Fprintf(p.Output, "%s ", title)

	slices.Sort(r.Hints)
	for _, hint := range r.Hints {
		if strings.HasSuffix(hint, "✔︎") {
			fmt.Fprintf(p.Output, "%s ", coloring.Green(hint))
		} else {
			fmt.Fprintf(p.Output, "%s ", coloring.Red(hint))
		}
	}

	fmt.Fprintf(
		p.Output,
		"for %s (%dx)",
		r.Duration.Round(time.Millisecond),
		r.Consecutive,
	)

	fmt.Fprint(p.Output, "\n")
}

// Summary outputs a total summary
func (p *Printer) Summary() {
	hosts := make([]string, 0, len(p.Streak))

	for h := range p.Streak {
		hosts = append(hosts, h)
	}

	sort.Strings(hosts)

	for _, host := range hosts {
		streak := p.Streak[host]

		if streak.LongestFailureStreak == 0 {
			fmt.Fprintf(p.Output, "- Longest outage %s: no outage\n", host)
		} else {
			fmt.Fprintf(
				p.Output,
				"- Longest outage %s: %s (%dx)\n",
				host,
				streak.LongestOutage.Round(time.Millisecond),
				streak.LongestFailureStreak,
			)
		}
	}
}
