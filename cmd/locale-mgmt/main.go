package main

import (
	"log"
	"os"

	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/handling"
	"github.com/ekr-paolo-carraro/locale-mgmt/pkg/session"
)

func main() {

	err := session.InitSessionStorage()
	if err != nil {
		log.Println(err.Error())
		return
	}

	r, err := handling.NewHandler()
	if err != nil {
		log.Fatalf("startup router give error:%s\n", err)
		return
	}

	log.Println(r.Run(":" + os.Getenv("PORT")))

}
