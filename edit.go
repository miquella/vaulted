package main

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/miquella/ask"
	"github.com/miquella/vaulted/lib"
	"golang.org/x/crypto/ssh"
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
		password, err = ask.HiddenAsk("New Password: ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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
	print("a - AWS Key")
	print("s - SSH Keys")
	print("v - Variables")
	print("? - Help")
	print("q - Quit")
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

func sshKeysHelp() {
	color.Set(color.FgYellow)
	print("")
	print("a - Add")
	print("d - Delete")
	print("? - Help")
	print("b - Back")
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

func edit(name string, v *vaulted.Vault) {
	exit := false
	for exit == false {
		cyan := color.New(color.FgCyan)
		cyan.Printf("\nVault: ")
		fmt.Printf("%s", name)
		printVariables(v)
		printAWS(v, false)
		printSSHKeys(v)

		input := readMenu("\nEdit vault: [a,s,v,?,q]: ")
		switch input {
		case "a":
			aws(v)
		case "s":
			sshKeysMenu(v)
		case "v":
			variables(v)
		case "q":
			exit = true
		case "?", "help":
			mainMenu()
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

func sshKeysMenu(v *vaulted.Vault) {
	exit := false

	for exit == false {
		printSSHKeys(v)
		input := readMenu("\nEdit ssh keys: [a,d,?,b]: ")
		switch input {
		case "a":
			addSSHKey(v)
		case "d":
			key := readValue("Key: ")
			_, ok := v.SSHKeys[key]
			if ok {
				delete(v.SSHKeys, key)
			} else {
				color.Red("Key '%s' not found", key)
			}
		case "b":
			exit = true
		case "?", "help":
			sshKeysHelp()
		default:
			color.Red("Command not recognized")
		}
	}
}

func addSSHKey(v *vaulted.Vault) {
	homeDir := ""
	user, err := user.Current()
	if err == nil {
		homeDir = user.HomeDir
	} else {
		homeDir = os.Getenv("HOME")
	}

	defaultFilename := ""
	filename := ""
	if homeDir != "" {
		defaultFilename = filepath.Join(homeDir, ".ssh", "id_rsa")
		filename = readValue(fmt.Sprintf("Key file (default: %s): ", defaultFilename))
		if filename == "" {
			filename = defaultFilename
		}
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(filepath.Join(homeDir, ".ssh"), filename)
		}
	} else {
		filename = readValue("Key file: ")
	}

	decryptedBlock, err := loadAndDecryptKey(filename)
	if err != nil {
		color.Red("%v", err)
		return
	}

	comment := loadPublicKeyComment(filename + ".pub")
	var name string
	if comment != "" {
		name = readValue(fmt.Sprintf("Name (default: %s): ", comment))
		if name == "" {
			name = comment
		}
	} else {
		name = readValue("Name: ")
		if name == "" {
			name = filename
		}
	}

	if v.SSHKeys == nil {
		v.SSHKeys = make(map[string]string)
	}
	v.SSHKeys[name] = string(pem.EncodeToMemory(decryptedBlock))
}

func loadAndDecryptKey(filename string) (*pem.Block, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	if x509.IsEncryptedPEMBlock(block) {
		var passphrase string
		var decryptedBytes []byte
		for i := 0; i < 3; i++ {
			passphrase, err = ask.HiddenAsk("Passphrase: ")
			if err != nil {
				return nil, err
			}

			decryptedBytes, err = x509.DecryptPEMBlock(block, []byte(passphrase))
			if err == nil {
				break
			}
			if err != x509.IncorrectPasswordError {
				return nil, err
			}
		}

		if err != nil {
			return nil, err
		}

		return &pem.Block{
			Type:  block.Type,
			Bytes: decryptedBytes,
		}, nil
	}
	return block, nil
}

func loadPublicKeyComment(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}

	_, comment, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return ""
	}
	return comment
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

func print(message string) {
	fmt.Printf("%s\n", message)
}

func printVariables(v *vaulted.Vault) {
	color.Cyan("\nVariables:")
	if len(v.Vars) > 0 {
		var keys []string
		for key, _ := range v.Vars {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			fmt.Printf("  %s: %s\n", key, v.Vars[key])
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

func printSSHKeys(v *vaulted.Vault) {
	color.Cyan("\nSSH Keys:")
	if len(v.SSHKeys) > 0 {
		keys := []string{}
		for key, _ := range v.SSHKeys {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			fmt.Printf("  %s\n", key)
		}
	} else {
		print("  [Empty]")
	}
}

func readMenu(message string) string {
	blue := color.New(color.FgBlue)
	blue.Printf(message)
	input := readInput(message)
	print("")
	return input
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
	return strings.TrimSpace(input)
}
