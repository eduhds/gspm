#!/bin/bash

set -e

appname=gspm
version="$(git describe --tags --abbrev=0)"
arch=x86_64
output=${appname}_${version}_Linux_${arch}

mkdir -p dist

rm -rf dist/AppDir 2> /dev/null || true
rm dist/*.AppImage 2> /dev/null || true
mkdir dist/AppDir

cmd="linuxdeploy \
    --appdir dist/AppDir \
    --executable dist/${appname}_linux_amd64*/$appname \
    --desktop-file linux/$appname.desktop \
     --output appimage \
"

sizes=(16 32 64 96 128 256)

for size in "${sizes[@]}"; do
    cmd="$cmd --icon-file linux/icons/${size}x${size}/$appname.png \
"
done

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-appimage \
    bash -c "$cmd"

mv *.AppImage dist/$output.AppImage
