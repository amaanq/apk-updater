#!/usr/bin/env bash

# The GPLv3 License (GPLv3)
#
# Copyright (c) 2023 Amaan Qureshi <amaanq12@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# NOTE this is a script to build for various platforms, 386 is the same as x86, amd64 is the same as x86_64

package=$1
if [[ -z "$package" ]]; then
	echo "usage: $0 <package-name>"
	exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "linux/amd64" "linux/386" "linux/arm64" "linux/arm" "darwin/amd64" "darwin/arm64")

for platform in "${platforms[@]}"; do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name='apk-updater-'$GOOS'-'$GOARCH
	if [ "$GOOS" = "windows" ]; then
		output_name+='.exe' # lol
	fi

	# want debuggable builds? build it yourself!
	GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_name $package
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi

	upx -9 "$output_name"
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done

zip -9 executables.zip apk-updater-*
