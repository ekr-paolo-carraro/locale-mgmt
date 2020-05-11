package authorizating

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/session"
	"github.com/gin-gonic/gin"
)

//LoginHandler manage login call to auth provider
func LoginHandler(c *gin.Context) {

	//random to generate state for request and then compare the code for getting auth-token
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//generate random state
	state := base64.StdEncoding.EncodeToString(b)

	//int session to store sate
	ss, err := session.Store.Get(c.Request, "auth-session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//store state and save session
	ss.Values["state"] = state
	err = ss.Save(c.Request, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//init connectoin with authenticator
	authenticator, err := NewAutenticator()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//call auth provider with random state
	redirectLocation := authenticator.Config.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, redirectLocation)
}

//LogoutHandler manage logout call
func LogoutHandler(c *gin.Context) {
	domain := os.Getenv("AUTH0_DOMAIN")
	logoutUrl, err := url.Parse(domain)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logoutUrl.Path += "v2/logout"
	params := url.Values{}

	var protocol string
	if c.Request.TLS == nil {
		protocol = "http"
	} else {
		protocol = "https"
	}

	urlToReturn := c.Request.Host
	if strings.Contains(urlToReturn, "http") == true {
		urlToReturn = strings.Replace(urlToReturn, "http://", "", 0)
		urlToReturn = strings.Replace(urlToReturn, "https://", "", 0)
	}
	log.Println(protocol + "://" + urlToReturn + "/welcome")
	returnTo, err := url.Parse(protocol + "://" + urlToReturn + "/welcome")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	params.Add("returnTo", returnTo.String())
	params.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = params.Encode()

	c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

//InfoHandler show version of server api and logged user info
func InfoHandler(c *gin.Context) {
	var msg map[string]interface{} = make(map[string]interface{})
	msg["version"] = "0.0.1"

	ss, err := session.Store.Get(c.Request, "auth-session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	candidate := ss.Values["profile"]
	if candidate != nil {
		msg["user"] = candidate
	}
	candidate = ss.Values["access_token"]
	if candidate != nil {
		msg["access_token"] = candidate
	}
	candidate = ss.Values["id_token"]
	if candidate != nil {
		msg["id_token"] = candidate
	}

	c.JSON(http.StatusOK, msg)
}
