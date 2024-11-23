package oidc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/un1uckyyy/go-telegram-oidc/internal/config"
	"golang.org/x/oauth2"
)

type Service interface {
	GetAuthUrl(userId int64) (string, error)
	CompleteAuth(ctx context.Context, state string, code string) (*oauth2.Token, error)
}

type TemporaryUserTracker interface {
	GetTmpUserInfo(state string) (telegramId int64, ok bool)
	SetTmpUserInfo(state string, telegramId int64) error
	PopTmpUserInfo(state string) (telegramId int64, ok bool)
}

type OidcService struct {
	Config      oauth2.Config
	UserTracker TemporaryUserTracker
}

func NewOidcService(tracker TemporaryUserTracker) (*OidcService, error) {
	endpoint := oauth2.Endpoint{
		AuthURL:  config.Instance.AuthUrl,
		TokenURL: config.Instance.TokenUrl,
	}
	redirectUrl, err := url.JoinPath(config.Instance.RedirectHost, "/auth")
	if err != nil {
		return nil, err
	}

	return &OidcService{
		Config: oauth2.Config{
			ClientID:     config.Instance.ClientId,
			ClientSecret: config.Instance.ClientSecret,
			RedirectURL:  redirectUrl,
			Endpoint:     endpoint,
		},
		UserTracker: tracker,
	}, nil
}

func (o *OidcService) GetAuthUrl(userId int64) string {
	state := uuid.NewString()
	_ = o.UserTracker.SetTmpUserInfo(state, userId)

	return o.Config.AuthCodeURL(state)
}

func (o *OidcService) CompleteAuth(ctx context.Context, state string, code string) (*oauth2.Token, error) {
	_, ok := o.UserTracker.GetTmpUserInfo(state)
	if !ok {
		return nil, fmt.Errorf("state not found")
	}

	token, err := o.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}
