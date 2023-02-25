/*
The GPLv3 License (GPLv3)

Copyright (c) 2023 Amaan Qureshi <amaanq12@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"sort"
	"strconv"
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

		game, err := selectGame("Which game do you want to download and decompress") // Have user pick a game
		if err != nil {
			return err
		}

		versions, err := apk.GetAllVersions(game.URL) // Get game versions
		if err != nil {
			return err
		}

		sort.SliceStable(versions, func(i, j int) bool { // Sort by order
			n_i_nums := strings.Split(versions[i].Version, ".")
			n_j_nums := strings.Split(versions[j].Version, ".")
			n_i, _ := strconv.Atoi(n_i_nums[0])
			n_j, _ := strconv.Atoi(n_j_nums[0])
			if n_i == n_j && len(n_i_nums) >= 2 && len(n_j_nums) >= 2 { // If both versions are of the same major version, but not the same minor version
				n_i, _ = strconv.Atoi(n_i_nums[1])
				n_j, _ = strconv.Atoi(n_j_nums[1])
				if n_i == n_j && len(n_i_nums) >= 3 && len(n_j_nums) >= 3 { // If both versions are of the same major version, minor version, but not build version
					n_i, _ = strconv.Atoi(n_i_nums[2])
					n_j, _ = strconv.Atoi(n_j_nums[2])
				}
			}
			return n_i > n_j
		})

		version, err := selectVersion(versions) // Have user pick a version
		if err != nil {
			return err
		}

		apk.Log.Infof("Downloading %s APK Version %s (Released on %s)\n", game.Name, version.Version, version.Date)
		_, err = apk.WgetAPK(game, version.DownloadURL, version.Version, outputDownloadFP) // Download the apk
		if err != nil {
			return err
		}
		apk.Log.Infof("Downloaded %s-%s.apk Successfully!", strings.ToLower(strings.ReplaceAll(game.Name, " ", "")), version.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&desiredVersion, "version", "v", "", "The desired version of the APK to download (default is newest, older versions unimplemented as of now)")
	downloadCmd.Flags().StringVarP(&outputDownloadFP, "output", "o", "", "Set the output folder for the decompressed APK (default is clash-major.minor.build")
}
