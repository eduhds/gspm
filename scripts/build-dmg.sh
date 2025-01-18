#!/bin/bash

set -e

appname=gspm
version="$(git describe --tags --abbrev=0)"
bundle=$appname.app
launch=launch
identifier=com.github.eduhds.$appname
output=${appname}_$version-macos-x86_64

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
    <string>${version:1}</string>
    <key>CFBundleVersion</key>
    <string>${version:1}</string>
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

cp dist/${appname}_darwin_amd64*/$appname dist/$bundle/Contents/MacOS
cp macos/$appname.icns dist/$bundle/Contents/Resources

npx appdmg macos/dmg.json dist/$output.dmg
