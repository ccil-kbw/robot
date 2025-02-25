#!/usr/bin/env bash

# Bump the version number in ./internal/version/version.go
# Usage: ./tools/bump_version.sh "v1.0.0"

set -e

CURRENT_VERSION=$(grep -oP 'const RobotVersion = "\K[^"]+' internal/version/version.go)
echo "Current version: $CURRENT_VERSION"

if [ -z "$1" ]; then
    echo "Usage: ./tools/bump_version.sh \"v1.0.0\""
    exit 1
fi

if [ "$CURRENT_VERSION" == "$1" ]; then
    echo "Version $1 is the same as the current version"
    exit 1
fi

echo "Bumping version to $1 in internal/version/version.go"
sed -i "s/const RobotVersion = \"$CURRENT_VERSION\"/const RobotVersion = \"$1\"/" internal/version/version.go

echo "Bumping version to $1 in README.md"
sed -i "s/# Masjid's Droid $CURRENT_VERSION/# Masjid's Droid $1/" README.md