package vaulted

import (
	"time"
)

type AWSKey struct {
	AWSCredentials
	MFA                     string        `json:"mfa,omitempty"`
	Role                    string        `json:"role,omitempty"`
	RoleDuration            time.Duration `json:"roleDuration,omitempty"`
	ForgoTempCredGeneration bool          `json:"forgoTempCredGeneration"`
}

func (k *AWSKey) Valid() bool {
	return k != nil && k.AWSCredentials.Valid()
}

func (k *AWSKey) RequiresMFA() bool {
	return k.Valid() && !k.ForgoTempCredGeneration && k.MFA != ""
}

func (k *AWSKey) GetAWSCredentials(duration time.Duration) (*AWSCredentials, error) {
	if k.ForgoTempCredGeneration {
		creds := k.AWSCredentials
		return &creds, nil
	}

	return k.AWSCredentials.GetSessionToken(duration)
}

func (k *AWSKey) GetAWSCredentialsWithMFA(mfaToken string, duration time.Duration) (*AWSCredentials, error) {
	if k.ForgoTempCredGeneration {
		creds := k.AWSCredentials
		return &creds, nil
	}

	return k.AWSCredentials.GetSessionTokenWithMFA(k.MFA, mfaToken, duration)
}
