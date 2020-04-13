package session

import (
	"encoding/gob"
	"os"

	"github.com/gorilla/sessions"
	"github.com/subosito/gotenv"
)

var (
	Store *sessions.FilesystemStore
)

//
func InitSessionStorage() error {

	gotenv.Load()

	Store = sessions.NewFilesystemStore("", []byte(os.Getenv("KEY_FOR_SESSION_STORE")))
	gob.Register(map[string]interface{}{})
	return nil
}
