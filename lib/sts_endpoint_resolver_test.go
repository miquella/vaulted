package vaulted_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/miquella/vaulted/lib"
)

var DefaultEndpoint = endpoints.ResolvedEndpoint{
	URL:                "https://sts.amazonaws.com",
	SigningRegion:      "us-east-1",
	SigningName:        "sts",
	SigningNameDerived: true,
	SigningMethod:      "v4",
}

type EndpointResolverTest struct {
	Service, Region string
	OptFns          []func(*endpoints.Options)

	ResolvedEndpoint endpoints.ResolvedEndpoint
}

func (t *EndpointResolverTest) HydratedEndpoint() endpoints.ResolvedEndpoint {
	e := DefaultEndpoint

	if t.ResolvedEndpoint.URL != "" {
		e.URL = t.ResolvedEndpoint.URL
	}

	if t.ResolvedEndpoint.SigningRegion != "" {
		e.SigningRegion = t.ResolvedEndpoint.SigningRegion
	}

	if t.ResolvedEndpoint.SigningName != "" {
		e.SigningName = t.ResolvedEndpoint.SigningName
	}

	if t.ResolvedEndpoint.SigningNameDerived {
		e.SigningNameDerived = t.ResolvedEndpoint.SigningNameDerived
	}

	if t.ResolvedEndpoint.SigningMethod != "" {
		e.SigningMethod = t.ResolvedEndpoint.SigningMethod
	}

	return e
}

func TestEndpointResolver(t *testing.T) {
	innerResolver := endpoints.ResolverFunc(func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		t.Fatal("Inner resolver called when it shouldn't be")
		return endpoints.ResolvedEndpoint{}, fmt.Errorf("Inner resolver called when it shouldn't be")
	})
	resolver := vaulted.STSEndpointResolver(innerResolver)

	var tests = []EndpointResolverTest{
		// Default global endpoint
		{},

		// Standard region
		{
			Region:           "us-west-2",
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "https://sts.us-west-2.amazonaws.com", SigningRegion: "us-west-2"},
		},

		// GovCloud region
		{
			Region:           "us-gov-west-1",
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "https://sts.us-gov-west-1.amazonaws.com", SigningRegion: "us-gov-west-1"},
		},

		// China region
		{
			Region:           "cn-north-1",
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "https://sts.cn-north-1.amazonaws.com.cn", SigningRegion: "cn-north-1"},
		},

		// Non-existent regions
		// (this is to ensure we can support regions that come online later without having to update the SDK to use them)
		{
			Region:           "us-nowhere-987",
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "https://sts.us-nowhere-987.amazonaws.com", SigningRegion: "us-nowhere-987"},
		},
		{
			Region:           "us-gov-nowhere-987",
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "https://sts.us-gov-nowhere-987.amazonaws.com", SigningRegion: "us-gov-nowhere-987"},
		},
		{
			Region:           "cn-nowhere-987",
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "https://sts.cn-nowhere-987.amazonaws.com.cn", SigningRegion: "cn-nowhere-987"},
		},

		// SSL disabled (why would you do this??)
		{
			OptFns:           []func(*endpoints.Options){func(options *endpoints.Options) { options.DisableSSL = true }},
			ResolvedEndpoint: endpoints.ResolvedEndpoint{URL: "http://sts.amazonaws.com"},
		},
	}

	for _, test := range tests {
		service := test.Service
		if service == "" {
			service = endpoints.StsServiceID
		}

		resolvedEndpoint, err := resolver.EndpointFor(service, test.Region, test.OptFns...)
		if err != nil {
			t.Errorf("Failed to resolve endpoint for service:%q region:%q: %v", service, test.Region, err)
		} else {
			expectedEndpoint := test.HydratedEndpoint()
			if !reflect.DeepEqual(resolvedEndpoint, expectedEndpoint) {
				t.Errorf("Incorrect endpoint returned for service:%q region:%q\n     Got: %+v\nExpected: %+v", service, test.Region, resolvedEndpoint, expectedEndpoint)
			}
		}
	}

	// Validate delegation for non-STS services
	delegated := false
	innerResolver = endpoints.ResolverFunc(func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		delegated = true
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	})
	resolver = vaulted.STSEndpointResolver(innerResolver)

	_, err := resolver.EndpointFor(endpoints.IamServiceID, "")
	if err != nil {
		t.Errorf("Failed to resolve endpoint for service:%q region:%q: %v", endpoints.IamServiceID, "", err)
	} else if !delegated {
		t.Error("Resolver did not delegate to next resolver as expected")
	}
}
