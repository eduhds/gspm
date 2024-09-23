#!/bin/bash

set -e

appname=gspm
version=0.1.0
os=$(go env GOOS)
arch=$(go env GOARCH)
output=$appname-$os-$arch

rm -rf $os/{AppDir,dist} 2> /dev/null
mkdir $os/{AppDir,dist}

cmd="LINUXDEPLOY_OUTPUT_APP_NAME=$output LINUXDEPLOY_OUTPUT_VERSION=$version linuxdeploy \
    --appdir $os/AppDir \
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

mv *.AppImage $os/dist/$appname-$os-$arch.AppImage

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-rpm \
    bash -c "LDNP_BUILD=rpm $cmd --output native_packages"

mv *.rpm $os/dist/$appname-$os-$arch.rpm

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-deb \
    bash -c "LDNP_BUILD=deb $cmd --output native_packages"

mv *.deb $os/dist/$appname-$os-$arch.deb

tar -C build/$os/$arch/release -czf $os/dist/$appname-$os-$arch.tar.gz $appname
