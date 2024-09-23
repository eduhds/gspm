#!/bin/bash

set -e

if [ "$(uname)" = "Darwin" ]; then
    if [ -f macos/gspm.icns ]; then
        exit 0
    fi
    
    mkdir -p macos/Icon.iconset
    
    sizes=(16 32 64 128 256 512)

    for size in "${sizes[@]}"; do
        sips -z $size $size assets/gspm-1024x1024.png --out macos/Icon.iconset/icon_${size}x${size}.png
        sips -z $((size*2)) $((size*2)) assets/gspm-1024x1024.png --out macos/Icon.iconset/icon_${size}x${size}@2x.png
    done

    iconutil -c icns -o macos/gspm.icns macos/Icon.iconset
else    
    mkdir -p linux/icons

    sizes=(16 32 64 96 128 256)

    for size in "${sizes[@]}"; do
        mkdir -p linux/icons/${size}x${size}

        if [ -f "linux/icons/${size}x${size}/gspm.png" ]; then
            continue
        fi

        ffmpeg -i assets/gspm-1024x1024.png -vf scale=${size}:${size} linux/icons/${size}x${size}/gspm.png
    done

    # Windows
    mkdir -p windows/icons

    sizes=(256)

    for size in "${sizes[@]}"; do
        if [ -f "windows/icons/icon-${size}.png" ]; then
            continue
        fi

        ffmpeg -i assets/gspm-1024x1024.png -vf scale=${size}:${size} windows/icons/icon-${size}.png
    done

    if ! [ -f "windows/icons/icon-256.png" ]; then
        ffmpeg -i windows/icons/icon-256.png windows/icons/icon.ico
    fi
fi
