/*
Copyright © 2022 Amaan Qureshi (aq0527@pm.me)

This file is part of the CLI application APK Updater.

This project, APK Updater, is not to be redistributed or copied without

the express permission of the copyright holder, Amaan Qureshi (amaanq).

*/
package cmd

import (
	"errors"
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
	Short: "Decompress an already existing apk or assets directory",
	Long: `This is to be used if you already have an APK downloaded, and want to simply decompress it. 

Set the APK file path using the -f flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputDecompressFP != "" && inputAssetsFP != "" {
			return errors.New("cannot specify both -f and -d")
		}

		if outputDecompressFP == "" {
			outputDecompressFP = "decompressed"
		}

		switch {
		case inputDecompressFP == "" && inputAssetsFP == "": // Default case
			game, err := selectGame() // Have user pick a game
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

			_bool := askToOnlyStoreAssets()
			apk.Log.Debug("Storing Assets:", _bool)

			apk.Log.Infof("Downloading %s APK Version %s (Released on %s)\n", game.Name, version.Version, version.Date)
			fp, err := apk.WgetAPK(game, version.DownloadURL, version.Version, "") // Download the apk to name-version.apk, return stored file path
			if err != nil {
				return err
			}

			apk.Log.Info("FP: ", fp)
			err = apk.DecompileAPK(fp) // Decompile this apk from file path above
			if err != nil {
				return err
			}

			// if err = apk.ExtractAssets("clash-" + version.Version); err != nil {
			// 	apk.Log.Error(err)
			// 	return err
			// }
			// apk.Log.Info("FP: ", fp)
			// inputAssetsFP = strings.TrimSuffix(fp, ".apk")
			// apk.Log.Info("Input Assets FP: ", inputAssetsFP, "output decompress FP: ", outputDecompressFP)
			// outputDecompressFP = "decompressed" + version.Version
			return apk.WalkAndDecompressAssets(".", outputDecompressFP)
		case inputDecompressFP != "":
			_, err := os.Stat(inputDecompressFP)
			if errors.Is(err, os.ErrNotExist) {
				return errors.New("Given APK File does not exist!")
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
			return apk.WalkAndDecompressAssets(inputAssetsFP, outputDecompressFP)
		case inputAssetsFP != "":
			_, err := os.Stat(inputAssetsFP)
			if errors.Is(err, os.ErrNotExist) {
				return errors.New("Given Assets Path does not exist!")
			}
			if err != nil {
				return err
			}
			return apk.WalkAndDecompressAssets(inputAssetsFP, outputDecompressFP)
		}
		return nil
	},
}

func selectGame() (*apk.GameLink, error) {
	templates := &promptui.SelectTemplates{
		Label: "		{{ . }}?",
		Active: "		     ↳ {{ .Name | cyan }}",
		Inactive: "			{{ .Name | cyan }}",
		Selected: "			{{ .Name | red | cyan }}",
		Details: `
			--------- Game ----------
			{{ "Name:" | faint }}	{{ .Name }}`,
	}
	prompt := promptui.Select{
		Label:     "Which game do you want to download and decompress",
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
		Label: "		{{ . }}?",
		Active: "		     ↳ {{ .Version | cyan }} ({{ .Date | red }})",
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
		Label:     "Do you want to only store the assets folder?",
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false
	}
	return result == "y" || result == "Y"
}

func init() {
	rootCmd.AddCommand(decompressCmd)
	decompressCmd.Flags().StringVarP(&inputDecompressFP, "file", "f", "", "Point to the APK to decompress")
	decompressCmd.Flags().StringVarP(&inputAssetsFP, "directory", "d", "", "Point to the assets folder to decompress")
	decompressCmd.Flags().StringVarP(&outputDecompressFP, "output", "o", "", "Set the output folder for the decompressed APK (default is clash-major.minor.build)")
}
