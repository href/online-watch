package online

import "fmt"

// Target points at a single host that needs to be checked for packetloss
type Target struct {

	// ID is the unique name or address of the target
	ID string

	// Label is an optional identifier
	Label string

	// Address points at the address of the host
	Host string

	// Check that determines whether the target is online
	Check Check
}

func (t *Target) Title() string {
	if t.Label == "" {
		return t.Host
	}

	return fmt.Sprintf("%s (%s)", t.Host, t.Label)
}
