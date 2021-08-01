#!/bin/bash

BUILD_VERSION=Development
BUILD_TYPE=Release
BUILD_CONFIG=amd64

usage()
{
    echo "Usage: $0"
    echo "  -c (--cfg) - indicates the build target such as amd64, armhf (pi 1, zero), etc."
    echo "  -t (--type) - indicates the build type, Release or Debug"
    echo "  -v (--version) - indicates the version for this build, default is Development"
    echo "  -h (--help) - print this usage statement"
}

while [ "$1" != "" ]; do
    case $1 in
        -t | --type )
            shift
            BUILD_TYPE=$1
            ;;
        -c | --cfg )
            shift
            BUILD_CONFIG=$1
            ;;
        -v | --version )
            shift
            BUILD_VERSION=$1
            ;;
        -h | --help )
            usage
            exit
            ;;
        * )
        usage
        exit 1
    esac
    shift
done

echo 'Linting...'
golangci-lint run

echo 'Building...'
go build -ldflags "-X site/cmd.Version=${BUILD_VERSION}" site.go
