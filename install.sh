curl -s https://github.com/labaneilers/drake/releases | grep releases.*windows-amd64 | sed -n -E 's/^.*download\/(.*)\/windows.*$/\1/p'