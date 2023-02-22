package main

import (
	"log"
	"server/internal/constants"
	"server/internal/user"
	"server/internal/ws"
	"server/pkg/db"
	"server/pkg/router"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalln("Error while creating db connection", err)
	}
	userRep := user.NewRepository(db.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start(constants.BASE_SERVER_URL)

}
