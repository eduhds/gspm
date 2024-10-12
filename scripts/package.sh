#!/bin/bash

set -e

appname=gspm
version=0.1.1
os=$(go env GOOS)
arch=$(go env GOARCH)
output=$appname-$os-$arch

mkdir -p dist

rm -rf dist/AppDir 2> /dev/null
rm dist/*.{rpm,deb,AppImage,tar.gz} 2> /dev/null

mkdir dist/AppDir

cmd="LINUXDEPLOY_OUTPUT_APP_NAME=$output LINUXDEPLOY_OUTPUT_VERSION=$version linuxdeploy \
    --appdir dist/AppDir \
    --executable build/$os/$arch/release/$appname \
    --desktop-file $os/$appname.desktop \
"

sizes=(16 32 64 96 128 256)

for size in "${sizes[@]}"; do
    cmd="$cmd --icon-file $os/icons/${size}x${size}/$appname.png \
"
done

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-appimage \
    bash -c "$cmd --output appimage"

mv *.AppImage dist/$output.AppImage

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-rpm \
    bash -c "LDNP_BUILD=rpm $cmd --output native_packages"

mv *.rpm dist/$output.rpm

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-deb \
    bash -c "LDNP_BUILD=deb $cmd --output native_packages"

mv *.deb dist/$output.deb

tar -C build/$os/$arch/release -czf dist/$output.tar.gz $appname
