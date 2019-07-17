package vaulted

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// The default endpoint resolver uses the global STS endpoint for all standard
// AWS regions, regardless of what region the client is configured to use. This
// resolver always uses the locally configured region instead.
func STSEndpointResolver(nextResolver endpoints.Resolver) endpoints.Resolver {
	return endpoints.ResolverFunc(func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.StsServiceID {
			var opts endpoints.Options
			opts.Set(optFns...)

			scheme := "https"
			if opts.DisableSSL {
				scheme = "http"
			}

			signingRegion := region
			hostname := ""
			if region == "" {
				signingRegion = "us-east-1"
				hostname = "sts.amazonaws.com"
			} else {
				p, found := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region)
				if !found {
					return endpoints.ResolvedEndpoint{}, fmt.Errorf("Unknown region: %s", region)
				}

				switch p.ID() {
				case "aws-cn":
					hostname = fmt.Sprintf("sts.%s.amazonaws.com.cn", region)
				case "aws", "aws-us-gov":
					fallthrough
				default:
					hostname = fmt.Sprintf("sts.%s.amazonaws.com", region)
				}
			}

			return endpoints.ResolvedEndpoint{
				URL:                fmt.Sprintf("%s://%s", scheme, hostname),
				SigningRegion:      signingRegion,
				SigningName:        service,
				SigningNameDerived: true,
				SigningMethod:      "v4",
			}, nil
		}

		return nextResolver.EndpointFor(service, region, optFns...)
	})
}
