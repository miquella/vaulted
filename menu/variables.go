package menu

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

type VariableMenu struct {
	*Menu
}

func (m *VariableMenu) Handler() error {
	var varErr error

	for {
		var input string
		m.Printer()
		if m.Vault.Vars == nil {
			input, varErr = interaction.ReadMenu("Edit environment variables: [a,b]: ")
		} else {
			input, varErr = interaction.ReadMenu("Edit environment variables: [a,S,D,b]: ")
		}
		if varErr != nil {
			return varErr
		}
		switch input {
		case "a", "add", "var", "variable", "variables":
			variableKey, keyErr := interaction.ReadValue("Name: ")
			if keyErr != nil {
				return keyErr
			}
			variableValue, valErr := interaction.ReadValue("Value: ")
			if valErr != nil {
				return valErr
			}
			if m.Vault.Vars == nil {
				m.Vault.Vars = make(map[string]string)
			}
			if _, exists := m.Vault.Vars[variableKey]; exists {
				confirm, err := interaction.ReadValue(fmt.Sprintf("Variable '%s' already exists. Overwrite? (y/n): ", variableKey))
				if err != nil {
					return err
				}
				if confirm != "y" {
					break
				}
			}
			m.Vault.Vars[variableKey] = variableValue
		case "S", "show", "hide":
			m.toggleHidden()
		case "D", "delete", "remove":
			variable, valErr := interaction.ReadValue("Variable name: ")
			if valErr != nil {
				return valErr
			}
			if _, exists := m.Vault.Vars[variable]; exists {
				delete(m.Vault.Vars, variable)
			} else {
				color.Red("Variable '%s' not found", variable)
			}
		case "b", "back":
			return nil
		case "q", "quit", "exit":
			var confirm string
			var err error
			confirm, err = interaction.ReadValue("Are you sure you wish to save and exit the vault? (y/n): ")
			if err == nil {
				if confirm == "y" {
					return ErrSaveAndExit
				}
			}
		case "?", "help":
			m.Help()
		default:
			color.Red("Command not recognized")
		}
	}
}

func (m *VariableMenu) Help() {
	color.Set(color.FgYellow)
	defer color.Unset()
	fmt.Println("")
	fmt.Println("a - Add")
	fmt.Println("S - Show/Hide Secrets")
	fmt.Println("D - Delete")
	fmt.Println("? - Help")
	fmt.Println("b - Back")
	fmt.Println("q - Quit")
}

func (m *VariableMenu) Printer() {
	color.Cyan("\nVariables:")
	if len(m.Vault.Vars) > 0 {
		var keys []string
		for key := range m.Vault.Vars {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			green.Printf("  %s: ", key)
			if m.Menu.ShowHidden {
				fmt.Printf("%s\n", m.Vault.Vars[key])
			} else {
				fmt.Printf("%s\n", faintColor.Sprint("<hidden>"))
			}
		}
	} else {
		fmt.Println("  [Empty]")
	}
}
