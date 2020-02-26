package menu

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/fatih/color"

	"github.com/miquella/vaulted/v3/lib"
)

type ImportCredentialsMenu struct {
	Menu
	Credentials *vaulted.AWSCredentials
}

func (m *ImportCredentialsMenu) Handler() error {
	credsChain := credentials.NewCredentials(&credentials.ChainProvider{
		Providers: []credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{},
		},
	})

	creds, err := credsChain.Get()
	if err != nil {
		return nil
	}

	if creds.SessionToken != "" {
		warningColor.Println("There appear to be AWS session credentials in your current environment.")
		warningColor.Println("Vaulted cannot import AWS session credentials.")
		return nil
	}

	for {
		warningColor.Println("There appear to be AWS credentials in your current environment.")
		input, err := interaction.ReadPrompt("Would you like to import these credentials? (Y/n): ")
		if err != nil {
			return err
		}

		switch strings.ToLower(input) {
		case "", "y", "yes":
			m.Credentials = &vaulted.AWSCredentials{
				ID:     creds.AccessKeyID,
				Secret: creds.SecretAccessKey,
			}
			return nil

		case "n", "no":
			return nil

		default:
			fmt.Println("")
			color.Red("Response not recognized. Please enter 'y' or 'n'.")
			fmt.Println("")
			continue
		}
	}
}
