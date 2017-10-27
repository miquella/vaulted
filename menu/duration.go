package menu

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/miquella/vaulted/lib"
)

type DurationMenu struct {
	Menu
}

func (m *DurationMenu) Handler() (err error) {
	var dur string
	maxDuration := 999 * time.Hour
	if m.Vault.AWSKey != nil && m.Vault.AWSKey.ForgoTempCredGeneration == false {
		maxDuration = 36 * time.Hour
	}
	readMessage := fmt.Sprintf("Duration (15m–%s): ", m.formatDuration(maxDuration))
	dur, err = interaction.ReadValue(readMessage)
	if err == nil {
		duration, durErr := time.ParseDuration(dur)
		if durErr != nil {
			color.Red("%s", durErr)
			return
		}
		if duration < 15*time.Minute || duration > maxDuration {
			errorMessage := fmt.Sprintf("Duration must be between 15m and %s", m.formatDuration(maxDuration))
			color.Red(errorMessage)
			return
		}
		m.Vault.Duration = duration
	}
	return
}

func (m *DurationMenu) formatDuration(duration time.Duration) string {
	dur := duration.String()
	if strings.HasSuffix(dur, "m0s") {
		dur = dur[:len(dur)-2]
	}
	if strings.HasSuffix(dur, "h0m") {
		dur = dur[:len(dur)-2]
	}
	return dur
}

func (m *DurationMenu) Printer() {
	cyan.Println("\nEnvironment:")
	green.Print("  Duration: ")
	var duration time.Duration
	if m.Vault.Duration == 0 {
		duration = vaulted.STSDurationDefault
	} else {
		duration = m.Vault.Duration
	}
	fmt.Printf("%s\n", m.formatDuration(duration))
}
