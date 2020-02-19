package menu

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fatih/color"
)

var (
	ErrInvalidCredentials = errors.New("Credentials are not valid")
	ErrNotAUser           = errors.New("Credentials are not for a user principal")
)

type DetectMFAMenu struct {
	*Menu
}

func (m *DetectMFAMenu) Handler() error {
	fmt.Println()
	fmt.Printf("Checking for MFA devices...\n")
	mfaDevices, err := m.getMFADevices()
	if err != nil {
		fmt.Printf("Unable to check for MFA devices.\n")
		fmt.Println()
		return err
	}

	switch len(mfaDevices) {
	case 0:
		fmt.Printf("No MFA devices found.\n")
		fmt.Println()

	case 1:
		fmt.Printf("1 MFA device found.\n")
		fmt.Println()
		return m.promptForSingle(mfaDevices[0])

	default:
		fmt.Printf("%d MFA devices found.\n", len(mfaDevices))
		fmt.Println()
		return m.promptForMultiple(mfaDevices)
	}

	return nil
}

func (m *DetectMFAMenu) promptForSingle(mfaDevice *iam.MFADevice) error {
	for {
		warningColor.Printf("  %s\n", *mfaDevice.SerialNumber)

		input, err := interaction.ReadPrompt("Would you like to use this MFA device? (Y/n): ")
		if err != nil {
			return err
		}

		switch strings.ToLower(input) {
		case "", "y", "yes":
			m.Vault.AWSKey.MFA = *mfaDevice.SerialNumber
			return nil

		case "n", "no":
			m.Vault.AWSKey.MFA = ""
			return nil

		default:
			fmt.Println()
			color.Red("Response not recognized. Please enter 'y' or 'n'.")
			fmt.Println()
		}
	}
}

func (m *DetectMFAMenu) promptForMultiple(mfaDevices []*iam.MFADevice) error {
	for {
		warningColor.Println("  0) None")
		for i, mfaDevice := range mfaDevices {
			warningColor.Printf("  %d) %s\n", i+1, *mfaDevice.SerialNumber)
		}

		input, err := interaction.ReadPrompt("Select the MFA device to use: ")
		if err != nil {
			return err
		}

		selection, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			fmt.Println()
			color.Red("Response not recognized. Please enter a value from 0 to %d.", len(mfaDevices))
			fmt.Println()
			continue
		}

		index := int(selection)
		if index < 0 || index > len(mfaDevices) {
			fmt.Println()
			color.Red("Invalid entry. Please enter a value from 0 to %d.", len(mfaDevices))
			fmt.Println()
			continue
		}

		if index == 0 {
			m.Vault.AWSKey.MFA = ""
		} else {
			m.Vault.AWSKey.MFA = *mfaDevices[index-1].SerialNumber
		}
		return nil
	}
}

func (m *DetectMFAMenu) iamClient() (*iam.IAM, error) {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			m.Vault.AWSKey.ID,
			m.Vault.AWSKey.Secret,
			m.Vault.AWSKey.Token,
		),
	}

	if m.Vault.AWSKey.Region != nil {
		config.Region = aws.String(*m.Vault.AWSKey.Region)
	}

	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return iam.New(s), nil
}

func (m *DetectMFAMenu) getMFADevices() ([]*iam.MFADevice, error) {
	// Must be valid credentials
	if !m.Vault.AWSKey.Valid() {
		return nil, ErrInvalidCredentials
	}

	// Use the caller identity to get the user's name
	callerIdentity, err := m.Vault.AWSKey.GetCallerIdentity()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(callerIdentity.Resource, "user/") {
		return nil, ErrNotAUser
	}

	username := path.Base(callerIdentity.Resource)

	// List the user's MFA devices
	iamClient, err := m.iamClient()
	if err != nil {
		return nil, err
	}

	var mfaDevices []*iam.MFADevice
	err = iamClient.ListMFADevicesPages(&iam.ListMFADevicesInput{UserName: &username}, func(output *iam.ListMFADevicesOutput, lastPage bool) bool {
		mfaDevices = append(mfaDevices, output.MFADevices...)

		return true
	})
	if err != nil {
		return nil, err
	}

	return mfaDevices, nil
}
