package handling

import (
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
		apiGroup.GET("/bundle/:bundleId/langs", authorizating.AuthRequired(), lph.GetAllLangs)

		apiGroup.GET("/locale-item/:id", authorizating.AuthRequired(), lph.GetLocaleItemById)
		apiGroup.POST("/locale-item", authorizating.AuthRequired(), lph.PostLocaleItem)
		apiGroup.POST("/locale-items", authorizating.AuthRequired(), lph.PostLocaleItems)

		apiGroup.POST("/locale-items/:bundle", authorizating.AuthRequired(), lph.GetLocaleItemByBundleKeyLang)

		apiGroup.DELETE("/locale-items/:bundle", authorizating.AuthRequired(), lph.DeleteLocaleItemByBundleKeyLang)
		apiGroup.DELETE("/locale-items/:bundle/lang/:langId", authorizating.AuthRequired(), lph.DeleteLocaleItemByBundleKeyLang)
		apiGroup.DELETE("/locale-items/:bundle/lang/:langId/key/:keyId", authorizating.AuthRequired(), lph.DeleteLocaleItemByBundleKeyLang)
		apiGroup.DELETE("/locale-items/:bundle/key/:keyId", authorizating.AuthRequired(), lph.DeleteLocaleItemByBundleKeyLang)

	}

	return rh, nil
}
