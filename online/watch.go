package online

import (
	"context"
	"fmt"
	"math"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

type Watch struct {
	Targets  []string
	NoTCP    bool
	NoICMP   bool
	Verbose  bool
	Timeout  int
	Interval int
	Port     int
}

// Run starts the workers watching the targets and returns once they are
// started. Results are then pushed out through the results channel.
//
// To stop the workers from running, cancel the given context.
//
// Returns an error if the targets could not be parsed
func (w *Watch) Run(ctx context.Context, results chan<- Result) error {
	targets, err := w.ParseTargets()
	if err != nil {
		return err
	}

	timeout := time.Duration(w.Timeout) * time.Millisecond

	for _, target := range targets {
		go func() {
			lastOkay := true
			lastChangeTime := time.Now()
			consecutive := 0
			last := time.Now()
			for {
				select {
				case <-time.After(max(timeout-time.Since(last), 0)):
					last = time.Now()
					result := target.Check.Execute(timeout)

					if result.Okay != lastOkay {
						lastOkay = result.Okay
						lastChangeTime = result.Done
						consecutive = 1
					} else {
						consecutive++
					}

					result.Consecutive = consecutive
					result.Duration = time.Since(lastChangeTime)
					result.Target = &target
					results <- result
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return nil
}

func (w *Watch) ParseTargets() ([]Target, error) {
	var err error

	if w.Port < 0 || w.Port > math.MaxUint16 {
		return nil, fmt.Errorf("not a valid port: %d", w.Port)
	}

	targets := make([]Target, 0, len(w.Targets))
	for _, target := range w.Targets {
		var label, host string

		// Extract the label first
		switch strings.Count(target, "=") {
		case 0:
			host = target
		case 1:
			label, host, _ = strings.Cut(target, "=")
		default:
			return nil, fmt.Errorf("not a valid target: %s", target)
		}

		// Parse the remainder for the host/port tuple
		host, port, err := w.ParseAddress(host)
		if err != nil {
			return nil, err
		}

		// Build the check
		check := CheckGroup{}

		if port == 0 {
			port = uint16(w.Port)
		}

		if !w.NoTCP {
			check.Checks = append(check.Checks, TCPCheck{
				Host: host,
				Port: port,
			})
		}

		if !w.NoICMP {
			check.Checks = append(check.Checks, ICMPCheck{
				Host: host,
			})
		}

		if len(check.Checks) == 0 {
			return nil, fmt.Errorf("no checks requested")
		}

		targets = append(targets, Target{
			ID:    host,
			Label: label,
			Host:  host,
			Check: &check,
		})
	}
	return targets, err
}

// ParseAddress takes a target without label, and returns the host and the
// port it includes, if given.
func (w *Watch) ParseAddress(host string) (string, uint16, error) {

	// Just an IP address
	addr, err := netip.ParseAddr(host)
	if err == nil {
		return addr.String(), 0, nil
	}

	// An IP address with port
	addrport, err := netip.ParseAddrPort(host)
	if err == nil {
		return addrport.Addr().String(), addrport.Port(), nil
	}

	// A hostname…
	switch strings.Count(host, ":") {
	case 0:
		// …without port
		return host, 0, nil
	case 1:
		// …with port
		host, strPort, _ := strings.Cut(host, ":")

		port, err := strconv.Atoi(strPort)
		if err != nil || port < 0 || port > math.MaxUint16 {
			return "", 0, fmt.Errorf("invalid port: %s", strPort)
		}

		return host, uint16(port), nil
	}

	return "", 0, fmt.Errorf("unable to parse host/port: %s", host)
}
