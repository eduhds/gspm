#!/bin/bash

set -e

appname=gspm
version=0.1.4
bundle=$appname.app
launch=launch
identifier=com.github.eduhds.$appname
output=$appname-macos-$(pkgx go env GOARCH)

mkdir -p dist

rm -rf dist/$bundle 2> /dev/null || true
rm dist/*.dmg 2> /dev/null || true

mkdir dist/$bundle
mkdir dist/$bundle/Contents
mkdir dist/$bundle/Contents/{MacOS,Resources}

cat << EOF > dist/$bundle/Contents/Info.plist 
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
    "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>English</string>
    <key>CFBundleExecutable</key>
    <string>$launch</string>
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

echo 'APPL?????' > dist/$bundle/Contents/PkgInfo

cat << EOF > dist/$bundle/Contents/MacOS/$launch
#!/bin/bash
EXECUTABLE=\$(dirname "\$0")/$appname
open "\$EXECUTABLE"
EOF

chmod +x dist/$bundle/Contents/MacOS/$launch

cp build/darwin/amd64/release/$appname dist/$bundle/Contents/MacOS
cp macos/$appname.icns dist/$bundle/Contents/Resources

npx appdmg macos/dmg.json dist/$output.dmg

cp LICENSE.txt build/darwin/amd64/release
cp README.md build/darwin/amd64/release
tar -C build/darwin/amd64/release -czf dist/$output.tar.gz $appname LICENSE.txt README.md
