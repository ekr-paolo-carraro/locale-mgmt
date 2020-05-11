package handling

import (
	"net/http"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/authorizating"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/storaging"
	"github.com/gin-gonic/gin"
)

type genericMessage struct {
	Message string
}

//NewHandler return a new router handler
func NewHandler() (*gin.Engine, error) {

	rh := gin.Default()

	rh.GET("/", welcomeHandler)
	rh.GET("/welcome", welcomeHandler)

	rh.GET("/callback", authorizating.CallbackHandler)
	rh.GET("/login", authorizating.LoginHandler)
	rh.GET("/logout", authorizating.LogoutHandler)

	rh.GET("/info", authorizating.InfoHandler)
	lph, err := storaging.NewPersistenceHandler()
	if err != nil {
		return nil, err
	}

	apiGroup := rh.Group("/api/v1")
	{
		apiGroup.GET("/restricted", authorizating.AuthRequired(), authorizating.RestrictedHandler)

		apiGroup.GET("/langs", authorizating.AuthRequired(), lph.GetAllLangs)
		apiGroup.GET("/bundles", authorizating.AuthRequired(), lph.GetAllBundles)

		apiGroup.POST("/locale-item", authorizating.AuthRequired(), lph.PostLocaleItemHandler)
		apiGroup.POST("/locale-items", authorizating.AuthRequired(), lph.PostLocaleItemsHandler)

		apiGroup.GET("/locale-items/:bundle", authorizating.AuthRequired(), lph.GetLocaleItemByBundleKeyLang)
		apiGroup.GET("/locale-items/:bundle/lang/:langId", authorizating.AuthRequired(), lph.GetLocaleItemByBundleKeyLang)
		apiGroup.GET("/locale-items/:bundle/lang/:langId/key/:keyId", authorizating.AuthRequired(), lph.GetLocaleItemByBundleKeyLang)
		apiGroup.GET("/locale-items/:bundle/key/:keyId", authorizating.AuthRequired(), lph.GetLocaleItemByBundleKeyLang)

	}

	return rh, nil
}

func welcomeHandler(c *gin.Context) {
	msg := genericMessage{}
	msg.Message = "Hi, welcome to locale-mgmt, don't know who you are so go to login"
	c.JSON(http.StatusOK, msg)
}
