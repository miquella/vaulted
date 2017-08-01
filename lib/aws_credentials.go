package vaulted

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	SIGNIN_URL       = "https://signin.aws.amazon.com/federation"
	ALLOW_ALL_POLICY = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "*",
				"Resource": "*"
			}
		]
	}`
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

func (c *AWSCredentials) GetFederationToken(duration time.Duration) (*AWSCredentials, error) {
	stsClient, err := c.stsClient()
	if err != nil {
		return nil, err
	}

	callerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	arn, err := ParseARN(*callerIdentity.Arn)
	if err != nil {
		return nil, err
	}

	getFederationToken, err := stsClient.GetFederationToken(&sts.GetFederationTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		Name:            aws.String(path.Base(arn.Resource)),
		Policy:          aws.String(ALLOW_ALL_POLICY),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(getFederationToken.Credentials), nil
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

func (c *AWSCredentials) AssumeRoleWithMFA(serialNumber, token, arn string, duration time.Duration) (*AWSCredentials, error) {
	stsClient, err := c.stsClient()
	if err != nil {
		return nil, err
	}

	assumeRole, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(roleSessionName(stsClient)),
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(token),
	})
	if err != nil {
		return nil, err
	}

	return AWSCredentialsFromSTSCredentials(assumeRole.Credentials), nil
}

func (c *AWSCredentials) WithLocalDefault() (*AWSCredentials, error) {
	if !c.Valid() {
		awsCreds := credentials.NewCredentials(&credentials.ChainProvider{
			Providers: []credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
			},
		})

		creds, err := awsCreds.Get()
		if err != nil {
			return nil, err
		}

		return &AWSCredentials{
			ID:     creds.AccessKeyID,
			Secret: creds.SecretAccessKey,
			Token:  creds.SessionToken,
		}, nil
	}

	return c, nil
}

func (c *AWSCredentials) WithDefault() (*AWSCredentials, error) {
	if !c.Valid() {
		creds, err := defaults.Get().Config.Credentials.Get()
		if err != nil {
			return nil, err
		}

		return &AWSCredentials{
			ID:     creds.AccessKeyID,
			Secret: creds.SecretAccessKey,
			Token:  creds.SessionToken,
		}, nil
	}

	return c, nil
}

func (c *AWSCredentials) GetSigninToken(duration *time.Duration) (string, error) {
	// Load default credentials, if necessary
	var sess struct {
		Id    string `json:"sessionId"`
		Key   string `json:"sessionKey"`
		Token string `json:"sessionToken"`
	}

	c, err := c.WithDefault()
	if err != nil {
		return "", err
	}
	sess.Id = c.ID
	sess.Key = c.Secret
	sess.Token = c.Token

	// Get signin token
	sessionJson, err := json.Marshal(sess)
	if err != nil {
		return "", err
	}

	getTokenQuery := make(url.Values)
	getTokenQuery.Set("Action", "getSigninToken")
	if duration != nil {
		getTokenQuery.Set("SessionDuration", strconv.Itoa(int(duration.Seconds())))
	}
	getTokenQuery.Set("Session", string(sessionJson))

	signinUrl, _ := url.Parse(SIGNIN_URL)
	signinUrl.RawQuery = getTokenQuery.Encode()
	resp, err := http.Get(signinUrl.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to retrieve federated signin token (status: %d)", resp.StatusCode)
	}

	var signinResponse struct {
		SigninToken string
	}
	err = json.NewDecoder(resp.Body).Decode(&signinResponse)
	if err != nil {
		return "", err
	}

	return signinResponse.SigninToken, nil
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
