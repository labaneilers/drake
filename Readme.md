## Installation

### Windows

Run this in powershell:

```
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/labaneilers/drake/master/Install.ps1'))
```

### Mac OSX

Run this in bash:

```
wget -O - https://raw.githubusercontent.com/labaneilers/drake/master/install-darwin-amd64.sh | bash
```