package handling

import (
	"net/http"
	"os"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/authorizating"
	"github.com/gin-gonic/gin"
)

//Handler router with persistence delegate
type Handler struct {
	PersistenceDelegate interface{}
}

//NewHandler return a new router handler
func NewHandler(delegate interface{}) error {

	rh := gin.Default()

	rh.GET("/", welcomeHandler)
	rh.GET("/welcome", welcomeHandler)

	rh.GET("/callback", authorizating.CallbackHandler)
	rh.GET("/login", authorizating.LoginHandler)
	rh.GET("/logout", authorizating.LogoutHandler)

	rh.GET("/info", authorizating.InfoHandler)

	apiGroup := rh.Group("/api/v1")
	{
		apiGroup.GET("/restricted", authorizating.AuthRequired(), authorizating.RestrictedHandler)
	}

	return rh.Run(os.Getenv("PORT"))
}

func welcomeHandler(c *gin.Context) {
	msg := genericMessage{}
	msg.Message = "Hi, welcome to locale-mgmt, don't know who you are so go to login"
	c.JSON(http.StatusOK, msg)
}

type genericMessage struct {
	Message string
}
