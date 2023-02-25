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
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/amaanq/apk-updater/apk"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var inputDecompressFP string
var inputAssetsFP string
var outputDecompressFP string

// decompressCmd represents the decompress command
var decompressCmd = &cobra.Command{
	Use:   "decompress",
	Short: "Decompress an APK",
	Long: `Decompress works in 3 ways:

1. Simply run decompress with no flags and follow the terminal prompts. The apk will be automatically downloaded and parsed for you.

2. Run decompress with the -f flag. This will decompress the APK file specified by the -f flag.

3. Run decompress with the -d flag. This will decompress the assets folder of an already DECOMPILED APK specified by the -d flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputDecompressFP != "" && inputAssetsFP != "" {
			return errors.New("cannot specify both -f and -d")
		}

		switch {
		case inputDecompressFP == "" && inputAssetsFP == "": // Default case
			game, err := selectGame("Which game do you want to download and decompress") // Have user pick a game
			if err != nil {
				return err
			}

			apk.Log.Info(game.URL)

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

			_bool := askToOnlyStoreAssets()

			_sc := askToDecompressDotSCFiles()
			if _sc {
				game.ValidDirectories = append(game.ValidDirectories, "sc")
			} else {
				apk.Log.Info("Not decompressing .sc files")
			}

			apk.Log.Infof("Downloading %s APK Version %s (Released on %s)\n", game.Name, version.Version, version.Date)
			fp, err := apk.WgetAPK(game, version.DownloadURL, version.Version, "") // Download the apk to name-version.apk, return stored file path .apk
			if err != nil {
				return err
			}

			err = apk.DecompileAPK(fp) // Decompile this apk from file path above (same path as apk without .apk)
			if err != nil {
				return err
			}

			if outputDecompressFP == "" {
				outputDecompressFP = defaultAssetOutputFolder(game, version)
			}

			if strings.TrimSuffix(fp, ".apk") == outputDecompressFP { // in case they're matching directories
				outputDecompressFP += "/decompressed"
			}
			assetsFP, err := apk.WalkAndDecompressAssets(game.ValidDirectories, strings.TrimSuffix(fp, ".apk"), outputDecompressFP)
			if err != nil {
				return err
			}

			if _bool {
				_ = apk.CleanUp(assetsFP, fp)
				assetsFP = "decompressed"
			}
			apk.Log.Infof("Done! Decompressed assets stored in ./%s\n", assetsFP)
		case inputDecompressFP != "":
			game, err := selectGame("What game is this (needed for knowing what folders to parse..)") // Have user pick a game
			if err != nil {
				return err
			}

			// query for bool
			_bool := askToOnlyStoreAssets()

			_sc := askToDecompressDotSCFiles()
			if _sc {
				game.ValidDirectories = append(game.ValidDirectories, "sc")
			}

			_, err = os.Stat(inputDecompressFP)
			if errors.Is(err, os.ErrNotExist) {
				return errors.New("given APK File does not exist")
			}
			if err != nil {
				return err
			}

			if !strings.HasSuffix(inputDecompressFP, ".apk") {
				return errors.New("invalid file path, must end in .apk")
			}
			err = apk.DecompileAPK(inputDecompressFP)
			if err != nil {
				return err
			}
			inputAssetsFP = strings.TrimSuffix(inputDecompressFP, ".apk")

			if outputDecompressFP == "" {
				outputDecompressFP = strings.TrimSuffix(inputDecompressFP, ".apk") + "-decompressed"
			}
			assetsFP, err := apk.WalkAndDecompressAssets(game.ValidDirectories, inputAssetsFP, outputDecompressFP)
			if err != nil {
				return err
			}
			apk.Log.Infof("Assets stored in %s\n", assetsFP)
			if _bool {
				_ = apk.CleanUp(assetsFP, inputDecompressFP)
			}
		case inputAssetsFP != "":
			game, err := selectGame("What game is this (needed for knowing what folders to parse..)") // Have user pick a game
			if err != nil {
				return err
			}

			_bool := askToOnlyStoreAssets()

			_sc := askToDecompressDotSCFiles()
			if _sc {
				game.ValidDirectories = append(game.ValidDirectories, "sc")
			}

			_, err = os.Stat(inputAssetsFP)
			if errors.Is(err, os.ErrNotExist) {
				return errors.New("given Assets Path does not exist")
			}
			if err != nil {
				return err
			}

			if outputDecompressFP == "" {
				outputDecompressFP = inputAssetsFP + "-decompressed"
			}
			assetsFP, err := apk.WalkAndDecompressAssets(game.ValidDirectories, inputAssetsFP, outputDecompressFP)
			if err != nil {
				return err
			}
			apk.Log.Infof("Assets stored in %s\n", assetsFP)
			if _bool {
				_ = apk.CleanUp(assetsFP, inputAssetsFP)
			}
		}
		return nil
	},
}

func selectGame(_prompt string) (*apk.GameLink, error) {
	templates := &promptui.SelectTemplates{
		Label:    "		{{ . }}?",
		Active:   "		     ↳ {{ .Name | cyan }}",
		Inactive: "			{{ .Name | cyan }}",
		Selected: "			{{ .Name | red | cyan }}",
		Details: `
			--------- Game ----------
			{{ "Name:" | faint }}	{{ .Name }}`,
	}
	prompt := promptui.Select{
		Label:     _prompt,
		Items:     apk.AllGameLinks,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &apk.AllGameLinks[index], nil
}

func selectVersion(versions []apk.VersionData) (*apk.VersionData, error) {
	templates := &promptui.SelectTemplates{
		Label:    "		{{ . }}?",
		Active:   "		     ↳ {{ .Version | cyan }} ({{ .Date | red }})",
		Inactive: "			{{ .Version | cyan }} ({{ .Date | red }})",
		Selected: "			{{ .Version | red | cyan }}",
		Details: `
			--------- APK ----------
			{{ "Version:" | faint }}	{{ .Version }}
			{{ "Release Date:" | faint }}	{{ .Date }}`,
	}
	prompt := promptui.Select{
		Label:     "Which APK version do you want to download and decompress",
		Items:     versions,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &versions[index], nil
}

func askToOnlyStoreAssets() bool {
	prompt := promptui.Prompt{
		Label:     "Do you want to clean up all files but the decompress ones?",
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false
	}
	return result == "y" || result == "Y"
}

func askToDecompressDotSCFiles() bool {
	prompt := promptui.Prompt{
		Label:     "Do you want to decompress .sc files NOTE: This uses a LOT of memory, >=8GB of RAM is recommended?",
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false
	}
	return result == "y" || result == "Y"
}

func defaultAssetOutputFolder(game *apk.GameLink, version *apk.VersionData) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(strings.ReplaceAll(game.Name, " ", "")), version.Version)
}

func init() {
	rootCmd.AddCommand(decompressCmd)
	decompressCmd.Flags().StringVarP(&inputDecompressFP, "file", "f", "", "Point to the APK to decompress")
	decompressCmd.Flags().StringVarP(&inputAssetsFP, "directory", "d", "", "Point to the assets folder to decompress")
	decompressCmd.Flags().StringVarP(&outputDecompressFP, "output", "o", "", "Set the output folder for the decompressed APK (default is clash-major.minor.build)")
}
