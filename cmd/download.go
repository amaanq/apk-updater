/*
Copyright © 2022 Amaan Qureshi (aq0527@pm.me)

This file is part of the CLI application APK Updater.

This project, APK Updater, is not to be redistributed or copied without

the express permission of the copyright holder, Amaan Qureshi (amaanq).

*/
package cmd

import (
	"strings"

	"github.com/amaanq/apk-updater/apk"
	"github.com/spf13/cobra"
)

var desiredVersion string
var outputDownloadFP string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the Clash of Clans apk",
	Long:  `This will eventually support multiple versions, but for now the latest one is the only one available.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputDownloadFP != "" && !strings.HasSuffix(outputDownloadFP, ".apk") {
			apk.Log.Warnf("The given output file path (%s) does not end in .apk, this can cause issues down the road...", outputDownloadFP)
		}

		url, err := apk.GetDownloadURL(apk.ClashofClans.URL)
		if err != nil {
			return err
		}

		if _, err := apk.WgetAPK(&apk.ClashofClans, url, desiredVersion, outputDownloadFP); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&desiredVersion, "version", "v", "", "The desired version of the APK to download (default is newest, older versions unimplemented as of now)")
	downloadCmd.Flags().StringVarP(&outputDownloadFP, "output", "o", "", "Set the output folder for the decompressed APK (default is clash-major.minor.build")
}
