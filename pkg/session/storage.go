package session

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/subosito/gotenv"
)

var (
	Store *sessions.CookieStore
)

//InitSessionStorage startup storage for authentication
func InitSessionStorage() error {

	err := gotenv.Load()
	if err != nil {
		return err
	}

	log.Println("test--" + os.Getenv("KEY_FOR_SESSION_STORE"))
	Store = sessions.NewCookieStore([]byte(os.Getenv("KEY_FOR_SESSION_STORE")))
	gob.Register(map[string]interface{}{})
	return nil
}
