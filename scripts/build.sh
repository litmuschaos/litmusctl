#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

platforms=(
"darwin/386"
"darwin/amd64"
"linux/386"
"linux/amd64"
"linux/arm"
"linux/arm64"
"windows/386"
"windows/amd64"
"windows/arm"
)

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    echo 'Building' $GOOS-$GOARCH
    output_name='litmusctl-'$GOOS'-'$GOARCH

    env GOOS=$GOOS GOARCH=$GOARCH go build -v -o platforms/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
