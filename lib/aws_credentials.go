package vaulted

import (
	"fmt"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AWSCredentials struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Token  string `json:"token,omitempty"`
}

func AWSCredentialsFromSTSCredentials(creds *sts.Credentials) *AWSCredentials {
	return &AWSCredentials{
		ID:     *creds.AccessKeyId,
		Secret: *creds.SecretAccessKey,
		Token:  *creds.SessionToken,
	}
}

func (c *AWSCredentials) Valid() bool {
	return c != nil && c.ID != "" && c.Secret != ""
}

func (c *AWSCredentials) ValidSession() bool {
	return c.Valid() && c.Token != ""
}

func (c *AWSCredentials) GetSessionToken(duration time.Duration) (*AWSCredentials, error) {
	stsClient, err := c.stsClient()
	if err != nil {
		return nil, err
	}

	getSessionToken, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(getSessionToken.Credentials), nil
}

func (c *AWSCredentials) GetSessionTokenWithMFA(serialNumber, token string, duration time.Duration) (*AWSCredentials, error) {
	stsClient, err := c.stsClient()
	if err != nil {
		return nil, err
	}

	getSessionToken, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(token),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(getSessionToken.Credentials), nil
}

func (c *AWSCredentials) AssumeRole(arn string, duration time.Duration) (*AWSCredentials, error) {
	stsClient, err := c.stsClient()
	if err != nil {
		return nil, err
	}

	assumeRole, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(roleSessionName(stsClient)),
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(assumeRole.Credentials), nil
}

func (c *AWSCredentials) stsClient() (*sts.STS, error) {
	// if c is nil, the default credential provider chain is used
	// (yes, I know this seems a little weird)
	config := &aws.Config{}
	if c != nil && c.ID != "" {
		config.Credentials = credentials.NewStaticCredentials(
			c.ID,
			c.Secret,
			c.Token,
		)
	}

	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sts.New(s), nil
}

func roleSessionName(stsClient *sts.STS) string {
	roleSessionName := DefaultSessionName

	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err == nil {
		arn, err := ParseARN(*callerIdentity.Arn)
		if err == nil {
			roleSessionName = fmt.Sprintf("%s@%s", path.Base(arn.Resource), arn.AccountId)
		}
	}

	return roleSessionName
}
