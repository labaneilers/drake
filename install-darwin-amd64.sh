set -e

version=$(curl -L -s https://github.com/labaneilers/drake/releases/latest | grep releases.*windows-amd64 | sed -n -E 's/^.*download\/(.*)\/windows.*$/\1/p')
fileurl="https://github.com/labaneilers/drake/releases/download/${version}/darwin-amd64.drk"

if [ "$version" = "" ]; then
    echo "ERROR: Couldn't get most recent version of drk"
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