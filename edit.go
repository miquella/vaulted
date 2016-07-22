package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/miquella/vaulted/lib"
)

func (cli VaultedCLI) Edit() {
	if len(cli) != 2 {
		fmt.Fprintln(os.Stderr, "You must specify a vault to edit")
		os.Exit(255)
	}

	var password string
	var vault *vaulted.Vault
	var err error

	if vaulted.VaultExists(cli[1]) {
		password, vault, err = openVault(cli[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		vault = &vaulted.Vault{}
	}

	edit(cli[1], vault)

	if password == "" {
		password = getPassword()
	}

	err = vaulted.SealVault(password, cli[1], vault)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainMenu() {
	color.Set(color.FgYellow)
	print("")
	print("v - Variables")
	print("a - AWS Key")
	print("? - Help")
	print("q - Quit")
	color.Unset()
}

func variableMenu() {
	color.Set(color.FgYellow)
	print("")
	print("a - Add")
	print("d - Delete")
	print("? - Help")
	print("b - Back")
	color.Unset()
}

func awsMenu() {
	color.Set(color.FgYellow)
	print("")
	print("k - Key")
	print("m - MFA")
	print("r - Role")
	print("s - Show Key")
	print("d - Delete")
	print("? - Help")
	print("b - Back")
	color.Unset()
}

func edit(name string, v *vaulted.Vault) {
	exit := false
	for exit == false {
		cyan := color.New(color.FgCyan)
		cyan.Printf("\nVault: ")
		fmt.Printf("%s", name)
		printVariables(v)
		printAWS(v, false)

		input := readMenu("\nEdit vault: [v,a,?,q]: ")
		switch input {
		case "v":
			variables(v)
		case "a":
			aws(v)
		case "q":
			exit = true
		case "?", "help":
			mainMenu()
		default:
			color.Red("Command not recognized")
		}
	}
}

func variables(v *vaulted.Vault) {
	exit := false

	for exit == false {
		printVariables(v)
		input := readMenu("\nEdit environment variables: [a,d,?,b]: ")
		switch input {
		case "a":
			variableKey := readValue("Name: ")
			variableValue := readValue("Value: ")
			if v.Vars == nil {
				v.Vars = make(map[string]string)
			}
			v.Vars[variableKey] = variableValue
		case "d":
			variable := readValue("Variable name: ")
			_, ok := v.Vars[variable]
			if ok {
				delete(v.Vars, variable)
			} else {
				color.Red("Variable '%s' not found", variable)
			}
		case "b":
			exit = true
		case "?", "help":
			variableMenu()
		default:
			color.Red("Command not recognized")
		}
	}
}

func aws(v *vaulted.Vault) {
	exit := false
	show := false

	for exit == false {
		var input string
		printAWS(v, show)
		if v.AWSKey == nil {
			input = readMenu("\nEdit AWS key [k,?,b]: ")
		} else {
			input = readMenu("\nEdit AWS key [k,m,r,s,d,?,b]: ")
		}

		switch input {
		case "k":
			awsAccesskey := readValue("Key ID: ")
			awsSecretkey := readValue("Secret: ")
			v.AWSKey = &vaulted.AWSKey{
				ID:     awsAccesskey,
				Secret: awsSecretkey,
				MFA:    "",
				Role:   "",
			}
		case "m":
			if v.AWSKey != nil {
				awsMfa := readValue("MFA ARN or serial number: ")
				v.AWSKey.MFA = awsMfa
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "r":
			if v.AWSKey != nil {
				awsRole := readValue("Role ARN: ")
				v.AWSKey.Role = awsRole
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "s":
			if v.AWSKey != nil {
				show = !show
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "d":
			if v.AWSKey != nil {
				removeKey := readValue("Delete your AWS key? (y/n): ")
				if removeKey == "y" {
					v.AWSKey = nil
				}
			} else {
				color.Red("Must associate an AWS key with the vault first")
			}
		case "b":
			exit = true
		case "?", "help":
			awsMenu()
		default:
			color.Red("Command not recognized")
		}
	}
}

func print(message string) {
	fmt.Printf("%s\n", message)
}

func printVariables(v *vaulted.Vault) {
	color.Cyan("\nVariables:")
	if len(v.Vars) > 0 {
		for key, value := range v.Vars {
			fmt.Printf("  %s: %s\n", key, value)
		}
	} else {
		print("  [Empty]")
	}
}

func printAWS(v *vaulted.Vault, show bool) {
	color.Cyan("\nAWS Key:")
	if v.AWSKey != nil {
		fmt.Printf("  Key ID: %s\n", v.AWSKey.ID)
		if !show {
			fmt.Printf("  Secret: %s\n", "<hidden>")
		} else {
			fmt.Printf("  Secret: %s\n", v.AWSKey.Secret)
		}
		if v.AWSKey.MFA != "" {
			fmt.Printf("  MFA: %s\n", v.AWSKey.MFA)
		}
		if v.AWSKey.Role != "" {
			fmt.Printf("  Role: %s\n", v.AWSKey.Role)
		}
	} else {
		print("  [Empty]")
	}
}

func readMenu(message string) string {
	blue := color.New(color.FgBlue)
	blue.Printf(message)
	return readInput(message)
}

func readValue(message string) string {
	green := color.New(color.FgGreen)
	green.Printf(message)
	return readInput(message)
}

func readInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.TrimSpace(input)
	print("")
	return input
}
