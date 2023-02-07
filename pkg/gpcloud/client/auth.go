package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v12"
)

const defaultKeycloakHostname = "https://auth.g-portal.com/auth"
const defaultKeycloakRealm = "master"

type GPCloudAuth struct {
	jwtToken      *gocloak.JWT
	gocloakClient *gocloak.GoCloak
	expires       time.Time
	authOpts      *AuthOptions
}

type AuthOptions struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	Realm        *string
	Hostname     *string
}

func (authOptions *AuthOptions) GetRealm() string {
	if authOptions.Realm == nil {
		return defaultKeycloakRealm
	}
	return *authOptions.Realm
}
func (authOptions *AuthOptions) GetHostname() string {
	if authOptions.Hostname == nil {
		return defaultKeycloakHostname
	}
	return *authOptions.Hostname
}

func (a *GPCloudAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if a.expires.Before(time.Now()) {
		if err := a.refresh(); err != nil {
			if err2 := a.login(); err2 != nil {
				return nil, err2
			}
			return nil, err
		}
	}
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", a.jwtToken.AccessToken),
	}, nil
}

func (a *GPCloudAuth) RequireTransportSecurity() bool {
	return true
}

func NewAuth(opts *AuthOptions) (*GPCloudAuth, error) {
	auth := &GPCloudAuth{
		authOpts: opts,
	}
	auth.gocloakClient = gocloak.NewClient(opts.GetHostname())
	if err := auth.login(); err != nil {
		return nil, err
	}
	return auth, nil
}

func (a *GPCloudAuth) login() error {
	token, err := a.gocloakClient.Login(context.Background(), a.authOpts.ClientID, a.authOpts.ClientSecret, a.authOpts.GetRealm(), a.authOpts.Username, a.authOpts.Password)
	if err != nil {
		return err
	}
	a.jwtToken = token
	a.expires = time.Now().Add(time.Duration(token.ExpiresIn-10) * time.Second)
	return nil
}

func (a *GPCloudAuth) refresh() error {
	token, err := a.gocloakClient.RefreshToken(context.Background(), a.jwtToken.RefreshToken, a.authOpts.ClientID, a.authOpts.ClientSecret, a.authOpts.GetRealm())
	if err != nil {
		return err
	}
	a.jwtToken = token
	a.expires = time.Now().Add(time.Duration(token.ExpiresIn-10) * time.Second)
	return nil
}
