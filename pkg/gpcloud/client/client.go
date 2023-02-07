package client

import (
	"crypto/tls"
	"fmt"
	"log"

	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/auth/v1/authv1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/cloud/v1/cloudv1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/metadata/v1/metadatav1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/network/v1/networkv1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/payment/v1/paymentv1grpc"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const DefaultEndpoint = "grpc.g-portal.cloud:443"

type Client struct {
	grpcClient *grpc.ClientConn
}

type EndpointOverrideOption string

// CloudClient Returns the CloudServiceClient
func (c *Client) CloudClient() cloudv1grpc.CloudServiceClient {
	return cloudv1grpc.NewCloudServiceClient(c.grpcClient)
}

// AuthClient Returns the CloudServiceClient
func (c *Client) AuthClient() authv1grpc.AuthServiceClient {
	return authv1grpc.NewAuthServiceClient(c.grpcClient)
}

// MetadataClient Returns the MetadataServiceClient
func (c *Client) MetadataClient() metadatav1grpc.MetadataServiceClient {
	return metadatav1grpc.NewMetadataServiceClient(c.grpcClient)
}

// NetworkClient Returns the NetworkServiceClient
func (c *Client) NetworkClient() networkv1grpc.NetworkServiceClient {
	return networkv1grpc.NewNetworkServiceClient(c.grpcClient)
}

// PaymentClient Returns the PaymentServiceClient
func (c *Client) PaymentClient() paymentv1grpc.PaymentServiceClient {
	return paymentv1grpc.NewPaymentServiceClient(c.grpcClient)
}

// NewClient Returns a new GRPC client
func NewClient(extraOptions ...interface{}) (*Client, error) {
	cl := &Client{}

	var options []grpc.DialOption
	// Certificate pinning
	options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(getTLSOptions())))

	// User Agent
	options = append(options, grpc.WithUserAgent(fmt.Sprintf("GPCloud Golang Client [%s]", Version)))

	endpoint := DefaultEndpoint
	authenticationDefined := false
	for _, option := range extraOptions {
		if opt, ok := option.(grpc.DialOption); ok {
			options = append(options, opt)
			continue
		}
		if opt, ok := option.(EndpointOverrideOption); ok {
			endpoint = string(opt)
			continue
		}
		if opt, ok := option.(AuthProviderOption); ok && !authenticationDefined {
			log.Printf("Using auth provider: %T", opt)
			options = append(options, grpc.WithPerRPCCredentials(&AuthOption{
				Provider: &opt,
			}))
			authenticationDefined = true
			continue
		}
	}

	clientConn, err := grpc.Dial(endpoint, options...)
	if err != nil {
		return nil, err
	}

	cl.grpcClient = clientConn
	return cl, nil
}

func getTLSOptions() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
}
