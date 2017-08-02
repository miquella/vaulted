package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/miquella/vaulted/lib"
)

type Proxy struct {
	Bind        string
	Region      string
	Service     string
	UpstreamURL string
}

func (p *Proxy) Credentials() *credentials.Credentials {
	return defaults.Get().Config.Credentials
}

func (p *Proxy) ReverseProxy() (*httputil.ReverseProxy, error) {
	target, err := url.Parse(p.UpstreamURL)
	if err != nil {
		return nil, err
	}

	signer := v4.NewSigner(p.Credentials())
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		// Rewrite the request
		req.Host = ""
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = strings.TrimSuffix(target.Path, "/") + "/" + strings.TrimPrefix(req.URL.Path, "/")
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		req.Header.Set("Connection", "close")
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		// Read the body
		var body io.ReadSeeker
		if req.Body != nil {
			defer req.Body.Close()
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				req.URL = nil
				return
			}

			body = bytes.NewReader(bodyBytes)
		}

		// Sign the request
		_, err = signer.Sign(req, body, p.Service, p.Region, time.Now().Add(-30*time.Second))
		if err != nil {
			req.URL = nil
		}
	}
	return &httputil.ReverseProxy{Director: director}, nil
}

func (p *Proxy) Run(store vaulted.Store) error {
	reverseProxy, err := p.ReverseProxy()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", p.Bind)
	if err != nil {
		return err
	}

	addr := listener.Addr().String()
	addr = strings.Replace(addr, "[::]", "[::1]", 1)
	addr = strings.Replace(addr, "0.0.0.0", "127.0.0.1", 1)
	fmt.Printf("Listening at http://%v\n", addr)

	return http.Serve(listener, reverseProxy)
}
