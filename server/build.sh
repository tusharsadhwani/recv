#!/usr/bin/env bash
set -euo pipefail

rm -rf dist
mkdir -p dist

function build_platform {
  if [[ -z $1 ]]; then
    echo 'Error: GOOS not provided'
    exit 1
  fi
  export GOOS=$1

  if [[ -z $2 ]]; then
    echo 'Error: GOARCH not provided'
    exit 1
  fi
  export GOARCH=$2

  case $GOARCH in
    amd64)
      PLATFORM='x86-64'
      ;;
    386)
      PLATFORM='x86'
      ;;
  esac

  FILETYPE="$GOOS-$PLATFORM"

  if [[ $GOOS = windows ]]; then
    go build -ldflags='-s' -o "dist/recv-$FILETYPE.exe" ./cmd/cli
    go build -ldflags='-s' -o "dist/recv-server-$FILETYPE.exe" ./cmd/server
  else
    go build -ldflags='-s' -o 'dist/recv' ./cmd/cli
    go build -ldflags='-s' -o 'dist/server' ./cmd/server
    tar -czf "dist/recv-$FILETYPE.tar.gz" dist/recv
    tar -czf "dist/recv-server-$FILETYPE.tar.gz" dist/server
    rm dist/recv dist/server
  fi
  
}

build_platform "linux" "386"
build_platform "linux" "amd64"
build_platform "darwin" "amd64"
build_platform "windows" "386"
build_platform "windows" "amd64"

echo "All files built and saved in ./dist"
