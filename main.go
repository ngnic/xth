package main

import "xendit-takehome/github/server"

func main() {
	app := server.NewApp()
	app.Run()
}
