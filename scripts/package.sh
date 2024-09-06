#!/bin/bash

set -e

version=0.0.4
output=gspm-$(go env GOOS)-$(go env GOARCH)

rm -rf build/AppDir 2> /dev/null
mkdir build/AppDir

cmd="LINUXDEPLOY_OUTPUT_APP_NAME=$output LINUXDEPLOY_OUTPUT_VERSION=$version linuxdeploy \
    --appdir build/AppDir \
    --executable build/gspm \
    --icon-file linux/gspm.png \
    --desktop-file linux/gspm.desktop \
"

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-appimage \
    bash -c "$cmd --output appimage"

mv *.AppImage build

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-rpm \
    bash -c "LDNP_BUILD=rpm $cmd --output native_packages"

mv *.rpm build

docker run --rm -v $(pwd):/builder eduhds/linuxdeploy-deb \
    bash -c "LDNP_BUILD=deb $cmd --output native_packages"

mv *.deb build
