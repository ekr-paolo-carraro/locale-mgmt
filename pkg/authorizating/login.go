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

//LoginHandler manage login call
func LoginHandler(c *gin.Context) {

	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)

	ss, err := session.Store.Get(c.Request, "auth-session")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ss.Values["state"] = state
	err = ss.Save(c.Request, c.Writer)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	authenticator, err := NewAutenticator()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authenticator.Config.AuthCodeURL(state))
}

//LogoutHandler manage logout call
func LogoutHandler(c *gin.Context) {
	domain := os.Getenv("OAUTH_PROVIDER")
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
	params.Add("client_id", os.Getenv("CLIENT_ID"))
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

	msg["user"] = ss.Values["profile"]
	c.JSON(http.StatusOK, msg)
}
