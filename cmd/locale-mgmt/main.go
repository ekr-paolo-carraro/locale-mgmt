package main

import (
	"log"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/handling"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/session"
)

func main() {

	err := session.InitSessionStorage()
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(handling.NewHandler(nil))

}
