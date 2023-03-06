package auth

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"time"
)

type ProviderKeycloakUserPassword struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	Hostname     *string
	Realm        *string

	tokenData *struct {
		Token   *gocloak.JWT
		Expires time.Time
	}
}

func (p *ProviderKeycloakUserPassword) refresh() error {
	hostname := hostnameDefault
	if p.Hostname != nil {
		hostname = *p.Hostname
	}
	realm := realmDefault
	if p.Realm != nil {
		realm = *p.Realm
	}
	if p.tokenData == nil {
		token, err := gocloak.NewClient(hostname).Login(context.Background(), p.ClientID, p.ClientSecret, realm, p.Username, p.Password)
		if err != nil {
			p.tokenData = nil
			return err
		}
		p.tokenData = &struct {
			Token   *gocloak.JWT
			Expires time.Time
		}{Token: token, Expires: time.Now().Add(time.Duration(token.ExpiresIn-10) * time.Second)}
	} else {
		token, err := gocloak.NewClient(hostname).RefreshToken(context.Background(), p.tokenData.Token.RefreshToken, p.ClientID, p.ClientSecret, realm)
		if err != nil {
			p.tokenData = nil
			return err
		}
		p.tokenData = &struct {
			Token   *gocloak.JWT
			Expires time.Time
		}{Token: token, Expires: time.Now().Add(time.Duration(token.ExpiresIn-10) * time.Second)}
	}
	return nil
}

func (p *ProviderKeycloakUserPassword) GetToken(ctx context.Context) (string, error) {
	if p.tokenData == nil || p.tokenData.Expires.Before(time.Now()) {
		if err := p.refresh(); err != nil {
			return "", err
		}
	}
	return p.tokenData.Token.AccessToken, nil
}
