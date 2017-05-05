package vaulted

import (
	"reflect"
	"testing"
)

func TestParseARN(t *testing.T) {
	testParseARN(t,
		"arn:aws:elasticbeanstalk:us-east-1:123456789012:environment/My App/MyEnvironment",
		&ARN{
			Partition: "aws",
			Service:   "elasticbeanstalk",
			Region:    "us-east-1",
			AccountId: "123456789012",
			Resource:  "environment/My App/MyEnvironment",
		},
	)

	testParseARN(t,
		"arn:aws:iam::123456789012:user/David",
		&ARN{
			Partition: "aws",
			Service:   "iam",
			Region:    "",
			AccountId: "123456789012",
			Resource:  "user/David",
		},
	)

	testParseARN(t,
		"arn:aws:rds:eu-west-1:123456789012:db:mysql-db",
		&ARN{
			Partition: "aws",
			Service:   "rds",
			Region:    "eu-west-1",
			AccountId: "123456789012",
			Resource:  "db:mysql-db",
		},
	)

	testParseARN(t,
		"arn:aws:s3:::my_corporate_bucket/exampleobject.png",
		&ARN{
			Partition: "aws",
			Service:   "s3",
			Region:    "",
			AccountId: "",
			Resource:  "my_corporate_bucket/exampleobject.png",
		},
	)
}

func TestParseARNFail(t *testing.T) {
	testParseARNFail(t, "blah")
	testParseARNFail(t, ":aws:iam::123456789012:user/David")
	testParseARNFail(t, "arn:aws:s3::")
}

func testParseARN(t *testing.T, arn string, expected *ARN) {
	actual, err := ParseARN(arn)
	if err != nil {
		t.Errorf("Failed to parse ARN: %s\n%v", arn, err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %#v\nGot: %#v", expected, actual)
		return
	}
}

func testParseARNFail(t *testing.T, arn string) {
	_, err := ParseARN(arn)
	if err == nil {
		t.Errorf("Expected to fail to parse ARN: %s", arn)
		return
	}
}
