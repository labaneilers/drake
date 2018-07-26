version=$(curl -s https://github.com/labaneilers/drake/releases | grep releases.*windows-amd64 | sed -n -E 's/^.*download\/(.*)\/windows.*$/\1/p')
fileurl="https://github.com/labaneilers/drake/releases/download/${version}/darwin-amd64.drk"

installdir="/usr/local/bin"
mkdir -p $installdir

filepath="${installdir}/drk"

echo $filepath
echo $fileurl

wget $fileurl -O $filepath
chmod 777 $filepath