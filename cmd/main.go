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
	defer myApp.Close()
	err = myApp.Run()
	if err != nil {
		log.Fatal(err)
	}

}
