package vaulted

import (
	"errors"
	"strings"
)

// http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
//   arn:partition:service:region:account-id:resource
//   arn:partition:service:region:account-id:resourcetype/resource
//   arn:partition:service:region:account-id:resourcetype:resource

var (
	ErrNotAnArn = errors.New("Not an ARN")
)

type ARN struct {
	Partition string
	Service   string
	Region    string
	AccountId string
	Resource  string
}

func ParseARN(arn string) (*ARN, error) {
	components := strings.SplitN(arn, ":", 6)
	if len(components) != 6 {
		return nil, ErrNotAnArn
	}

	if components[0] != "arn" {
		return nil, ErrNotAnArn
	}

	return &ARN{
		Partition: components[1],
		Service:   components[2],
		Region:    components[3],
		AccountId: components[4],
		Resource:  components[5],
	}, nil
}

func (a *ARN) String() string {
	components := []string{
		"arn",
		a.Partition,
		a.Service,
		a.Region,
		a.AccountId,
		a.Resource,
	}

	if a.Partition == "" {
		components[1] = "aws"
	}

	return strings.Join(components, ":")
}
