package menu

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

type Interaction struct {
	rlValue  *readline.Instance
	rlOutput *readline.Instance
}

func (interaction *Interaction) ReadMenu(message string) (string, error) {
	if interaction.rlOutput == nil {
		var err error
		interaction.rlOutput, err = readline.New("")
		if err != nil {
			return "", err
		}
	}

	fmt.Println("")
	menuColor.Println("?=Help; q=Save and Quit; Ctrl+C=Abort")
	input, err := interaction.ReadInput(menuColor.Sprint(message), interaction.rlOutput)
	fmt.Println("")
	return input, err
}

func (interaction *Interaction) ReadValue(message string) (string, error) {
	if interaction.rlValue == nil {
		var err error
		interaction.rlValue, err = readline.New("")
		if err != nil {
			return "", err
		}
	}
	return interaction.ReadInput(color.GreenString(message), interaction.rlValue)
}

func (interaction *Interaction) ReadPrompt(message string) (string, error) {
	if interaction.rlValue == nil {
		var err error
		interaction.rlValue, err = readline.New("")
		if err != nil {
			return "", err
		}
	}
	return interaction.ReadInput(warningColor.Sprint(message), interaction.rlValue)
}

func (interaction *Interaction) ReadInput(message string, rl *readline.Instance) (string, error) {
	rl.SetPrompt(message)
	line, err := rl.Readline()
	if err == readline.ErrInterrupt {
		return "", ErrAbort
	}
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
