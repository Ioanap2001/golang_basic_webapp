package main

import (
	"log"
)

func main() {
	app, err := initializeApp()
	if err != nil {
		log.Fatal(err)
	}

	srv := app.getServer()

	app.infoLog.Printf("Starting server on %s", app.serverCfg.Address)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
