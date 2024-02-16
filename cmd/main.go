package main

import (
	conf "example1/config"
	"example1/internal/app"
)

func main() {
	config := conf.GetConfig()

	myApp := app.New(config)

	myApp.Run()

}
