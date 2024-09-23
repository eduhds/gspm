#!/bin/bash

set -e

appname=gspm
version=0.1.0
bundle=dist/$appname.app
executable=launch
identifier=com.github.eduhds.$appname

mkdir -p macos
mkdir -p macos/dist

rm -rf macos/dist/* 2> /dev/null

mkdir macos/$bundle
mkdir macos/$bundle/Contents
mkdir macos/$bundle/Contents/{MacOS,Resources}

cat << EOF > macos/$bundle/Contents/Info.plist 
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
    "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>English</string>
    <key>CFBundleExecutable</key>
    <string>$executable</string>
    <key>CFBundleIdentifier</key>
    <string>$identifier</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>$appname</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>$version</string>
    <key>CFBundleVersion</key>
    <string>$version</string>
    <key>CFBundleIconFile</key>
    <string>$appname.icns</string>
    <key>LSUIElement</key>
    <true/>
    <key>CFBundleURLTypes</key>
    <array>
      <dict>
        <key>CFBundleURLName</key>
        <string>$identifier</string>
        <key>CFBundleURLSchemes</key>
        <array>
          <string>$appname</string>
        </array>
      </dict>
    </array>
    <key>NSAppleScriptEnabled</key>
    <true/>
  </dict>
</plist>
EOF

echo 'APPL?????' > macos/$bundle/Contents/PkgInfo

cat << EOF > macos/$bundle/Contents/MacOS/$executable
#!/bin/bash
EXECUTABLE=\$(dirname "\$0")/$appname
open -a Terminal --args "\$EXECUTABLE" "\$@"
EOF

chmod +x macos/$bundle/Contents/MacOS/$executable

cp build/darwin/amd64/release/$appname macos/$bundle/Contents/MacOS
cp macos/$appname.icns macos/$bundle/Contents/Resources

dmg_file=macos/dist/$appname-$(pkgx go env GOOS)-$(pkgx go env GOARCH).dmg

npx appdmg macos/dmg.json $dmg_file

