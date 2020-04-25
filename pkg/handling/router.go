package handling

import (
	"net/http"
	"os"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/authorizating"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/storaging"
	"github.com/gin-gonic/gin"
)

type genericMessage struct {
	Message string
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
	lph, err := storaging.NewPersistenceHandler()
	if err != nil {
		return err
	}

	apiGroup := rh.Group("/api/v1")
	{
		apiGroup.GET("/restricted", authorizating.AuthRequired(), authorizating.RestrictedHandler)
		apiGroup.POST("/locale-item", authorizating.AuthRequired(), lph.PostLocaleItemHandler)
	}

	return rh.Run(":" + os.Getenv("PORT"))
}

func welcomeHandler(c *gin.Context) {
	msg := genericMessage{}
	msg.Message = "Hi, welcome to locale-mgmt, don't know who you are so go to login"
	c.JSON(http.StatusOK, msg)
}
