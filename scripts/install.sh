#!/bin/bash

set -e

COLLOR_OFF='\033[0m'
COLLOR_RED='\033[0;31m'
COLLOR_GREEN='\033[0;32m'
COLLOR_BLUE='\033[0;34m'

printf "$COLLOR_BLUE"

cat << 'EOF'

 ,adPPYb,d8  ,adPPYba,  8b,dPPYba,   88,dPYba,,adPYba,   
a8"    `Y88  I8[    ""  88P'    "8a  88P'   "88"    "8a  
8b       88   `"Y8ba,   88       d8  88      88      88  
"8a,   ,d88  aa    ]8I  88b,   ,a8"  88      88      88  
 `"YbbdP"Y8  `"YbbdP"'  88`YbbdP"'   88      88      88  
 aa,    ,88             88                               
  "Y8bbdP"              88                               

EOF

printf "$COLLOR_OFF"

version=v0.2.1
os=$(uname -s)
arch=amd64
url=https://github.com/eduhds/gspm/releases/download/$version/gspm-${os,}-$arch.tar.gz

echo "📦 Installing gspm..."
echo ""

if command -v curl &> /dev/null; then
    curl -L $url -o /tmp/gspm.tar.gz
elif command -v wget &> /dev/null; then
    wget -O /tmp/gspm.tar.gz $url
else
    echo "❌ Please install curl or wget and try again."
    exit 1
fi

echo ""

tar -C /tmp -xzf /tmp/gspm.tar.gz && \
    sudo mv /tmp/gspm /usr/local/bin && \
    sudo chmod +x /usr/local/bin/gspm

echo ""

if [ $? -ne 0 ]; then
    echo "❌ Failed to install gspm."
    echo "See https://github.com/eduhds/gspm for more information."
    exit 1
else
    echo '✅ gspm installed successfully!'
fi
