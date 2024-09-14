#!/bin/bash

set -e

appname=gspm
version=0.0.4
os=$(go env GOOS)
arch=$(go env GOARCH)
output=$appname-$os-$arch

rm -rf $os/{AppDir,dist} 2> /dev/null
mkdir $os/{AppDir,dist}

cmd="LINUXDEPLOY_OUTPUT_APP_NAME=$output LINUXDEPLOY_OUTPUT_VERSION=$version linuxdeploy \
    --appdir $os/AppDir \
    --executable build/$os/$arch/$appname \
    --icon-file $os/$appname.png \
    --desktop-file $os/$appname.desktop \
"

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-appimage \
    bash -c "$cmd --output appimage"

mv *.AppImage $os/dist

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-rpm \
    bash -c "LDNP_BUILD=rpm $cmd --output native_packages"

mv *.rpm $os/dist

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-deb \
    bash -c "LDNP_BUILD=deb $cmd --output native_packages"

mv *.deb $os/dist

tar -C build/$os/$arch -czf $os/dist/$appname-$os-$arch.tar.gz $appname
