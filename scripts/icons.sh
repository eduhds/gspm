#!/bin/bash

set -e

if [ "$(uname)" = "Darwin" ]; then
    if [ -f macos/gspm.icns ]; then
        exit 0
    fi
    
    mkdir -p macos/Icon.iconset
    
    sizes=(16 32 64 128 256 512)

    for size in "${sizes[@]}"; do
        sips -z $size $size gspm-1024x1024.png --out macos/Icon.iconset/icon_${size}x${size}.png
        sips -z $((size*2)) $((size*2)) gspm-1024x1024.png --out macos/Icon.iconset/icon_${size}x${size}@2x.png
    done

    iconutil -c icns -o macos/gspm.icns macos/Icon.iconset
else
    if [ -f linux/icons/**/gspm.png ]; then
        exit 0
    fi
    
    mkdir -p linux/icons

    sizes=(16 32 64 96 128 256)

    for size in "${sizes[@]}"; do
        mkdir -p linux/icons/${size}x${size}
        ffmpeg -i gspm-1024x1024.png -vf scale=${size}:${size} linux/icons/${size}x${size}/gspm.png
    done
fi
