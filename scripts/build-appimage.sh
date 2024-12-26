#!/bin/bash

set -e

appname=gspm
version="$(git describe --tags --abbrev=0)"
os=linux
arch=x86_64
output=${appname}_$version-$arch

mkdir -p dist

rm -rf dist/AppDir 2> /dev/null || true
mkdir dist/AppDir

cmd="linuxdeploy \
    --appdir dist/AppDir \
    --executable dist/${appname}_${os}_amd64*/$appname \
    --desktop-file $os/$appname.desktop \
     --output appimage \
"

sizes=(16 32 64 96 128 256)

for size in "${sizes[@]}"; do
    cmd="$cmd --icon-file $os/icons/${size}x${size}/$appname.png \
"
done

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-appimage \
    bash -c "$cmd"

mv *.AppImage dist/$output.AppImage
