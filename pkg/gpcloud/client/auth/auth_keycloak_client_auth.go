/* This authentication method is not fully implemented (yet)
   but will be used as long life access token replacement later on */

package auth

import (
	"context"
	"time"

	"github.com/Nerzal/gocloak/v12"
)

type ProviderKeycloakClientAuth struct {
	ClientID     string
	ClientSecret string
	Hostname     *string
	Realm        *string

	tokenData *struct {
		Token   *gocloak.JWT
		Expires time.Time
	}
}

func (p *ProviderKeycloakClientAuth) refresh() error {
	hostname := hostnameDefault
	if p.Hostname != nil {
		hostname = *p.Hostname
	}
	realm := realmDefault
	if p.Realm != nil {
		realm = *p.Realm
	}
	token, err := gocloak.NewClient(hostname).GetToken(context.Background(), realm, gocloak.TokenOptions{
		ClientID:     &p.ClientID,
		ClientSecret: &p.ClientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
	})
	if err != nil {
		p.tokenData = nil
		return err
	}
	p.tokenData = &struct {
		Token   *gocloak.JWT
		Expires time.Time
	}{Token: token, Expires: time.Now().Add(time.Duration(token.ExpiresIn-10) * time.Second)}
	return nil
}

func (p *ProviderKeycloakClientAuth) GetToken(ctx context.Context) (string, error) {
	if p.tokenData == nil || p.tokenData.Expires.Before(time.Now()) {
		if err := p.refresh(); err != nil {
			return "", err
		}
	}
	return p.tokenData.Token.AccessToken, nil
}
