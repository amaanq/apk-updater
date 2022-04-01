#!/usr/bin/env bash

# Copyright (C) 2022 Amaan Qureshi (aq0527@pm.me)
#
#
# This file is part of APK Updater.
#
#
# This project, APK Updater, is not to be redistributed or copied without
# 
# the express permission of the copyright holder, Amaan Qureshi (amaanq).


# NOTE this is a script to build for various platforms, 386 is the same as x86, amd64 is the same as x86_64

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "linux/amd64" "linux/386" "linux/arm64" "linux/arm" "darwin/amd64" "darwin/arm64")

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name='apk-updater-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe' # lol 
	fi	

	# want debuggable builds? build it yourself! 
	GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_name $package
    if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi

	upx -9 $output_name
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
zip -9 executables.zip apk-updater-*