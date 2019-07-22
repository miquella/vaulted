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
	durationMenu := &DurationMenu{Menu: &m.Menu}
	awsMenu := &AWSMenu{Menu: &m.Menu}
	variableMenu := &VariableMenu{Menu: &m.Menu}
	sshKeysMenu := &SSHKeyMenu{Menu: &m.Menu}

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
			err = durationMenu.Handler()
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

	fmt.Println("a,aws      - AWS Key")
	fmt.Println("s,ssh      - SSH Keys")
	fmt.Println("v,vars     - Variables")
	fmt.Println("d,duration - Session Duration")
	fmt.Println("S,show     - Show/Hide Secrets")
	fmt.Println("?,help     - Help")
	fmt.Println("q,quit     - Quit")
}
