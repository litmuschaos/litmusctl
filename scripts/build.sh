#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

platforms=(
"darwin/386"
"dragonfly/amd64"
"freebsd/386"
"freebsd/amd64"
"freebsd/arm"
"linux/386"
"linux/amd64"
"linux/arm"
"linux/arm64"
"netbsd/386"
"netbsd/amd64"
"netbsd/arm"
"openbsd/386"
"openbsd/amd64"
"openbsd/arm"
"plan9/386"
"plan9/amd64"
"solaris/amd64"
"windows/amd64"
"windows/386"
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
