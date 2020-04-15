package handling

import (
	"net/http"
	"os"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/authorizating"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/session"
	"github.com/gin-gonic/gin"
)

//Handler router with persistence delegate
type Handler struct {
	PersistenceDelegate interface{}
}

//NewHandler returnn a new router handler
func NewHandler(delegate interface{}) error {

	rh := gin.Default()

	rh.GET("/welcome", welcomeHandler)
	rh.GET("/", welcomeHandler)

	rh.GET("/callback", authorizating.CallbackHandler)
	rh.GET("/login", authorizating.LoginHandler)
	rh.GET("/logout", authorizating.LogoutHandler)
	rh.GET("/info", infoHandler)

	apiGroup := rh.Group("/api/v1")
	{
		apiGroup.GET("/test", authorizating.AuthRequired(), testApiHandler)
	}

	return rh.Run(os.Getenv("PORT"))
}

func infoHandler(c *gin.Context) {
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

func welcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, server is working: don't know who you are so go to login")
}

func testApiHandler(c *gin.Context) {
	c.String(http.StatusOK, "u r in protected area")
}
