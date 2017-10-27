package menu

import (
	"errors"

	"github.com/fatih/color"
	"github.com/miquella/vaulted/lib"
)

var (
	green = color.New(color.FgGreen)
	cyan  = color.New(color.FgCyan)
	blue  = color.New(color.FgBlue)

	faintColor   = color.New(color.Faint)
	menuColor    = color.New(color.FgHiBlue)
	warningColor = color.New(color.FgHiYellow)

	ErrAbort       = errors.New("Aborted by user. Vault unchanged.")
	ErrSaveAndExit = errors.New("Exiting at user request.")
)

type handler func() error
type output func()

var interaction = &Interaction{}

// Menu the type of all menus for the edit classes, provides a standardized interface for abstraction
type Menu struct {
	Vault      *vaulted.Vault
	ShowHidden bool
}

func (m *Menu) toggleHidden() {
	m.ShowHidden = !m.ShowHidden
}
