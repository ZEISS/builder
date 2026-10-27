package defaults

import (
	"math/rand"

	tea "charm.land/bubbletea/v2"
)

// CoalescedWheelMsg is emitted by the input filter after collapsing many
// terminal wheel samples into one message. DeltaX and DeltaY carry the summed
// scroll magnitude so handlers can scroll proportionally.
type CoalescedWheelMsg struct {
	Mouse  tea.Mouse
	DeltaX float64
	DeltaY float64
}

// Randomizer is an interface for generating random strings.
type Randomizer []string

// Random returns a random string from the randomizer.
func (r Randomizer) Random() string {
	if len(r) == 0 {
		return ""
	}
	return r[rand.Intn(len(r))]
}

// ReadyPlaceholders is a randomizer for ready placeholders.
var ReadyPlaceholders = Randomizer{
	"Ready!",
	"Ready...",
	"Ready?",
	"Ready for instructions",
}

// WorkingPlaceholders is a randomizer for working placeholders.
var WorkingPlaceholders = Randomizer{
	"Working!",
	"Working...",
	"Brrrrr...",
	"Prrrrrrrr...",
	"Processing...",
	"Thinking...",
}
