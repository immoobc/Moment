package main

import "moment/core"

func main() {
	if !core.EnsureSingleInstance() {
		return // another instance is already running
	}
	app := NewMomentApp()
	app.Run()
}
