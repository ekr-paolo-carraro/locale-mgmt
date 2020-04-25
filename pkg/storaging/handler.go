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
func (lph LocalePersistenceHandler) PostLocaleItemHandler(c *gin.Context) {
	var localeItem LocaleItem
	err := c.ShouldBind(&localeItem)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on bind payload: %v", err)}
		c.JSON(http.StatusBadRequest, msg)
	}

	localeItemReturned, err := lph.PersistenceDelegate.PostLocaleItem(localeItem)
	if err != nil {
		msg := ErrorMessage{fmt.Sprintf("Error on persist item: %v", err)}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	c.JSON(http.StatusOK, localeItemReturned)
}
