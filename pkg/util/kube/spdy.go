// Copied from https://github.com/istio/istio/blob/master/pkg/kube/spdy.go

package kube

import (
	"crypto/tls"
	"fmt"
	"net/http"

	spdyStream "k8s.io/apimachinery/pkg/util/httpstream/spdy"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
)

// roundTripperFor creates a SPDY upgrader that will work over custom transports.
func roundTripperFor(restConfig *rest.Config) (http.RoundTripper, spdy.Upgrader, error) {
	// Get the TLS config.
	tlsConfig, err := rest.TLSConfigFor(restConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting TLS config: %w", err)
	}
	if tlsConfig == nil && restConfig.Transport != nil {
		// If using a custom transport, skip server verification on the upgrade.
		// nolint: gosec
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	var upgrader *spdyStream.SpdyRoundTripper
	if restConfig.Proxy != nil {
		upgrader = spdyStream.NewRoundTripperWithProxy(tlsConfig, restConfig.Proxy)
	} else {
		upgrader = spdyStream.NewRoundTripper(tlsConfig)
	}
	wrapper, err := rest.HTTPWrappersForConfig(restConfig, upgrader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed creating SPDY upgrade wrapper: %w", err)
	}
	return wrapper, upgrader, nil
}
