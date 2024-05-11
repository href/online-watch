package online

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Check describes the interface of anything that can checked to be online
type Check interface {
	Execute(timeout time.Duration) Result
}

// Group combines multiple checks together and stores the results
type CheckGroup struct {
	Checks  []Check
	Results []Result
}

// Execute runs all checks concurrently and groups the result:
//
// - Result.Done is set to the end of the last check.
// - Result.Took is set to the longest check time.
// - Result.Okay is true if all checks passed.
func (g *CheckGroup) Execute(timeout time.Duration) Result {
	var wg sync.WaitGroup
	var mu sync.Mutex

	if g.Results != nil {
		g.Results = g.Results[:0]
	}

	for _, check := range g.Checks {
		wg.Add(1)
		go func(check Check) {
			defer wg.Done()

			r := check.Execute(timeout)
			mu.Lock()
			g.Results = append(g.Results, r)
			mu.Unlock()
		}(check)
	}

	wg.Wait()

	result := Result{Okay: true}
	for _, r := range g.Results {
		if r.Done.After(result.Done) {
			result.Done = r.Done
		}

		if !r.Okay {
			result.Okay = false
		}

		result.Hints = append(result.Hints, r.Hints...)
		result.Took = max(result.Took, r.Took)
	}

	return result
}

// Checks the availabilty of a host using `ping`
type ICMPCheck struct {
	Host string
}

func (c ICMPCheck) Execute(timeout time.Duration) Result {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := Result{}

	var exe string
	if strings.Contains(c.Host, ":") {
		exe = "ping6"
	} else {
		exe = "ping"
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, exe, "-c", "1", c.Host)
	err := cmd.Run()

	result.Done = time.Now()
	result.Took = time.Since(start)
	result.Okay = err == nil

	if exe == "ping6" {
		result.Hints = []string{
			fmt.Sprintf("ICMPv6 %s", result.ShortText())}
	} else {

		result.Hints = []string{
			fmt.Sprintf("ICMP %s", result.ShortText())}
	}

	return result
}

// Checks if a TCP handshake can be made
type TCPCheck struct {
	Host string
	Port uint16
}

func (c TCPCheck) Execute(timeout time.Duration) Result {
	socket := net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))

	result := Result{}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", socket, timeout)

	result.Done = time.Now()
	result.Took = time.Since(start)
	result.Okay = err == nil
	result.Hints = []string{
		fmt.Sprintf("TCP/%d %s", c.Port, result.ShortText())}

	if err == nil {
		conn.Close()
	}

	return result
}
