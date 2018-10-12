package menu

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// MainMenu is the top level menu entrypoint handler
type MainMenu struct {
	Menu
	VaultName string
}

func (m *MainMenu) Handler() error {
	var err error
	durationMenu := &DurationMenu{Menu: Menu{Vault: m.Vault}}
	awsMenu := &AWSMenu{Menu: Menu{Vault: m.Vault}, ShowHiddenVars: &m.ShowHidden}
	variableMenu := &VariableMenu{Menu: Menu{Vault: m.Vault}, ShowHiddenVars: &m.ShowHidden}
	sshKeysMenu := &SSHKeyMenu{Menu: Menu{Vault: m.Vault}}

	for {
		cyan.Printf("\nVault: ")
		fmt.Printf("%s", m.VaultName)
		variableMenu.Printer()
		awsMenu.Printer()
		sshKeysMenu.Printer()
		durationMenu.Printer()

		var input string
		input, err = interaction.ReadMenu("Edit vault: [a,s,v,d,S]: ")
		if err != nil {
			break
		}
		switch input {
		case "a", "aws":
			err = awsMenu.Handler()
		case "s", "ssh":
			err = sshKeysMenu.Handler()
		case "v", "vars", "variables":
			err = variableMenu.Handler()
		case "d", "duration":
			return durationMenu.Handler()
		case "S", "show", "hide":
			m.toggleHidden()
		case "b", "q", "quit", "exit":
			return nil
		case "?", "help":
			m.Help()
		default:
			color.Red("Command not recognized")
		}

		if err != nil {
			break
		}
	}

	if err == io.EOF || err == ErrSaveAndExit {
		return nil
	}
	return err
}

func (m *MainMenu) Help() {
	menuColor.Set()
	defer color.Unset()

	fmt.Println("a - AWS Key")
	fmt.Println("s - SSH Keys")
	fmt.Println("v - Variables")
	fmt.Println("d - Environment Duration")
	fmt.Println("S - Show/Hide Secrets")
	fmt.Println("? - Help")
	fmt.Println("q - Quit")
}
