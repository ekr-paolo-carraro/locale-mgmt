package storaging

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//LocalePersistenceHandler manages route for persistence
type LocalePersistenceHandler struct {
	PersistenceDelegate LocalePersistencer
}

//NewPersistenceHandler handles persitence request
func NewPersistenceHandler() (*LocalePersistenceHandler, error) {
	lph := &LocalePersistenceHandler{}

	lp, err := NewPostgresPersistenceService()
	if err != nil {
		return nil, err
	}

	lph.PersistenceDelegate = *lp

	return lph, nil
}

//PostLocaleItemHandler handle persitensce of a single locale item
func (lph LocalePersistenceHandler) PostLocaleItem(c *gin.Context) {
	var localeItem LocaleItem
	err := c.ShouldBind(&localeItem)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on bind payload: %v", err)}
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	if localeItem.isValid() == false {
		msg := ErrorMessage{fmt.Sprintf("Localeitem not valid: %v", localeItem)}
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	localeItemReturned, err := lph.PersistenceDelegate.PostLocaleItem(localeItem)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on persist item: %v", err)}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	c.JSON(http.StatusCreated, localeItemReturned)
}

//PostLocaleItemHandler handle persitensce of an array locale items
func (lph LocalePersistenceHandler) PostLocaleItems(c *gin.Context) {
	var localeItems []LocaleItem
	err := c.ShouldBind(&localeItems)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on bind payload: %v", err)}
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	numInserted, err := lph.PersistenceDelegate.PostLocaleItems(localeItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := MassiveResult{}
	result.NumSuccessfull = numInserted
	result.NumFailed = int64(len(localeItems)) - numInserted
	c.JSON(http.StatusCreated, result)
}

//GetLocaleItemHandler handle retrive for locale items
func (lph LocalePersistenceHandler) GetLocaleItemByBundleKeyLang(c *gin.Context) {
	var localeItems []LocaleItem
	pKey := c.Param("keyId")
	pLang := c.Param("langId")
	pBundle := c.Param("bundle")

	localeItems, err := lph.PersistenceDelegate.GetLocaleItems(pKey, pBundle, pLang)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on retrive items for %s, %s, %s : %v", pKey, pBundle, pLang, err)}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	c.JSON(http.StatusOK, localeItems)
}

//DeleteLocaleItemHandler handle retrive for delete locale items
func (lph LocalePersistenceHandler) DeleteLocaleItemByBundleKeyLang(c *gin.Context) {
	pKey := c.Param("keyId")
	pLang := c.Param("langId")
	pBundle := c.Param("bundle")

	numDeleteItems, err := lph.PersistenceDelegate.DeleteLocaleItems(pKey, pBundle, pLang)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on delete items for %s, %s, %s : %v", pKey, pBundle, pLang, err)}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	result := MassiveResult{}
	result.NumSuccessfull = numDeleteItems
	result.NumFailed = 0

	c.JSON(http.StatusOK, result)
}

//GetAllLangs return all lang
func (lph LocalePersistenceHandler) GetAllLangs(c *gin.Context) {
	result, err := lph.PersistenceDelegate.GetLangs()
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on retrive langs: %v", err)}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	c.JSON(http.StatusOK, result)
}

//GetAllLangs return all bundles
func (lph LocalePersistenceHandler) GetAllBundles(c *gin.Context) {
	result, err := lph.PersistenceDelegate.GetBundles()
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on retrive bundles: %v", err)}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	c.JSON(http.StatusOK, result)
}
