package client

import (
	"crypto/tls"
	"fmt"

	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/auth/v1/authv1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/cloud/v1/cloudv1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/metadata/v1/metadatav1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/network/v1/networkv1grpc"
	"buf.build/gen/go/gportal/gportal-cloud/grpc/go/gpcloud/api/payment/v1/paymentv1grpc"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const DefaultEndpoint = "grpc.gpcloud.space:443"

type Client struct {
	grpcClient *grpc.ClientConn
}

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
func NewClient(authOptions AuthOptions, options ...grpc.DialOption) (*Client, error) {
	cl := &Client{}

	// Certificate pinning
	options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(getTLSOptions())))

	// User Agent
	options = append(options, grpc.WithUserAgent(fmt.Sprintf("GPCloud Golang Client [%s]", Version)))

	auth, err := NewAuth(&authOptions)
	if err != nil {
		return nil, err
	}
	// Access Token
	options = append(options, grpc.WithPerRPCCredentials(auth))

	clientConn, err := grpc.Dial(DefaultEndpoint, options...)
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
