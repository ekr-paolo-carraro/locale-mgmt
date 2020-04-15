package authorizating

import (
	"context"
	"errors"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/session"
	"github.com/gin-gonic/gin"
)

type Autenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Context  context.Context
}

//NewAutenticator return an object with provider and other info to get auth token
func NewAutenticator() (*Autenticator, error) {

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, os.Getenv("OAUTH_PROVIDER"))
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
		Endpoint:     provider.Endpoint(),
	}

	return &Autenticator{
		Provider: provider,
		Config:   conf,
		Context:  ctx,
	}, nil
}

//CallbackHandler manage callback call
func CallbackHandler(c *gin.Context) {

	ss, err := session.Store.Get(c.Request, "auth-session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if c.Request.URL.Query().Get("state") != ss.Values["state"] {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Invalid state parameter"))
		return
	}

	authenticator, err := NewAutenticator()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), c.Request.URL.Query().Get("code"))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("No id_token field in oauth2 token"))
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("CLIENT_ID"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIdToken)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to verify id token: "+err.Error()))
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to marshall token: "+err.Error()))
		return
	}

	ss.Values["id_token"] = rawIdToken
	ss.Values["access_token"] = token.AccessToken
	ss.Values["profile"] = profile
	err = ss.Save(c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to save session: "+err.Error()))
		return
	}

	c.Redirect(http.StatusSeeOther, "/version")
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss, err := session.Store.Get(c.Request, "auth-session")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if _, ok := ss.Values["profile"]; !ok {
			c.Redirect(http.StatusSeeOther, "/welcome")
			return
		}

		c.Next()
	}
}
