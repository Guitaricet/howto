#!/bin/bash

# Determine the platform
if [[ "$(uname)" == "Darwin" ]]; then
    PLATFORM="darwin"
elif [[ "$(uname)" == "Linux" ]]; then
    PLATFORM="linux"
elif [[ "$(uname)" =~ "MINGW" ]]; then
    PLATFORM="windows"
    echo "Sorry, we don't have pre-built binaries for Windows yet, but it's very simple to build howto from source!"
    echo "git clone https://github.com/Guitaricet/howto"
    echo "cd howto"
    echo "go build"
    echo "If you don't have Go installed, follow https://go.dev for instrucitons"
    exit 1
else
    echo "Unsupported platform: $(uname)"
    exit 1
fi

# Determine the architecture
if [[ "$(uname -m)" == "x86_64" ]]; then
    ARCH="386"
elif [[ "$(uname -m)" == "arm64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $(uname -m)"
    exit 1
fi

# Query GitHub API for the latest release
LATEST_RELEASE=$(curl --silent "https://api.github.com/repos/Guitaricet/howto/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Download and install the correct binary
if [[ $PLATFORM == "windows" ]]; then
    URL="https://github.com/Guitaricet/howto/releases/download/$LATEST_RELEASE/howto-$LATEST_RELEASE-$PLATFORM-$ARCH.zip"
    curl -L $URL -o howto.zip
    unzip howto.zip
else
    URL="https://github.com/Guitaricet/howto/releases/download/$LATEST_RELEASE/howto-$LATEST_RELEASE-$PLATFORM-$ARCH.tar.gz"
    curl -L $URL | tar xz
fi

# Print success message
echo "Downloaded howto $LATEST_RELEASE to $(pwd)"
echo "Disclaimer: Howto suggestions are generated by an AI model and are not guaranteed to be safe to execute or to be executable at all. Please use common sense when using the suggested commands."
