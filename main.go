// Copyright (C) 2022 Amaan Qureshi (aq0527@pm.me)
//
//
// This file is part of APK Updater.
//
//
// This project, APK Updater, is not to be redistributed or copied without
//
// the express permission of the copyright holder, Amaan Qureshi (amaanq).

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
