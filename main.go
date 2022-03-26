package main

import (
	"github.com/amaanq/apk-updater/apk"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		apk.Log.Fatal(err)
	}

	if err := apk.UpdateAPK(); err != nil {
		apk.Log.Fatal(err)
	}
}
