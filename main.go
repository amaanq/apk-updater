package main

import "github.com/amaanq/apk-updater/apk"

func main() {
	if err := apk.CheckForNewVersion(); err != nil {
		panic(err)
	}
}
