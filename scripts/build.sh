#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

tag=$2

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

rm -rf platforms-$tag/*
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    echo 'Building' $GOOS-$GOARCH
    output_name='litmusctl-'$GOOS-$GOARCH

    env GOOS=$GOOS GOARCH=$GOARCH VERSION=$tag go build -v -o platforms-$tag/$output_name $package

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    cd platforms-$tag
    mv $output_name litmusctl
    tar -czvf $output_name-$tag.tar.gz litmusctl
    rm -rf litmusctl
    cd ..
done