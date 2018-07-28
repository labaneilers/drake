#!/usr/bin/env bash
set -e

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     os=linux;;
    Darwin*)    os=darwin;;
    CYGWIN*)    os=windows;;
    MINGW*)     os=windows;;
    *)          os="unknown"
esac

if [ "$os" = "unknown" ]; then
    echo "Unknown OS: $unameOut"
fi

unameArchOut="$(uname -m)"
case "${unameArchOut}" in
    x86_64*)    arch=amd64;;
    armv8*)     arch=arm;;
    *)          arch="unknown"
esac

if [ "$arch" = "unknown" ]; then
    echo "Unknown CPU: $unameArchOut"
fi

echo "Getting drk ($os $arch)..."

version=$(curl -L -s https://github.com/labaneilers/drake/releases/latest | grep releases.*windows-amd64 | sed -n 's/^.*download\/\(.*\)\/windows.*$/\1/p')
fileurl="https://github.com/labaneilers/drake/releases/download/${version}/${os}-${arch}.drk"

if [ "$version" = "" ]; then
    echo "ERROR: Couldn't get most recent version of drk for $os and $arch"
    exit 1
fi

installdir="/usr/local/bin"
mkdir -p $installdir

filepath="${installdir}/drk"
versionpath="${filepath}.version"

echo "Latest version: ${version}"

if [ -f $filepath ] && [ -f $versionpath ]; then
    diskVersion=$(cat $versionpath)
    if [ "$version" = "$diskVersion" ]; then
        echo "Already up-to-date"
        exit
    fi
fi

echo "Downloading from ${fileurl}..."

wget --quiet $fileurl -O $filepath
chmod 777 $filepath

echo $version > $versionpath

"drk successfully installed in $installdir"