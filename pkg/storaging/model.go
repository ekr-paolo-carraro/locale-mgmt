package storaging

import (
	"time"
)

//LocaleItem rappresents the item used for rappresent content in UI for every locale
type LocaleItem struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Bundle  string `json:"bundle"`
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

//LocaleItemHistory rappresents history traking for locale items
type LocaleItemHistory struct {
	ID               string
	LocaleItemID     string
	User             string
	ModificationDate time.Time
}

//ErrorMessage rappresents error message
type ErrorMessage struct {
	Message string
}
