package main

import (
	"Umbrella/internal/app"
	"log"
)

func main() {
	myApp, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = myApp.Run()
	if err != nil {
		log.Fatal(err)
	}
}
