if [ "$(uname)" = "Darwin" ]; then
    echo "macos"
else
    mkdir -p linux/icons

    sizes=(16 32 64 96 128 256)

    for size in "${sizes[@]}"; do
        mkdir -p linux/icons/${size}x${size}
        ffmpeg -i gspm-1024x1024.png -vf scale=${size}:${size} linux/icons/${size}x${size}/gspm.png
    done
fi
