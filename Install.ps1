$ErrorActionPreference = "Stop"

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$installDir = Join-Path $env:ALLUSERSPROFILE "drake"
if (![System.IO.Directory]::Exists($installDir)) {[void][System.IO.Directory]::CreateDirectory($installDir)}
$filePath = Join-Path $installDir "drk.exe"
$versionFilePath = "$($filePath).version"

function Get-LatestVersion() {
    $r = (Invoke-WebRequest "https://github.com/labaneilers/drake/releases").Content
    $line = $r.Split("`n") | Where-Object { $_ -match "download.*windows-amd64" }
    $found = $line -match 'download\/(.*)\/windows-amd64'
    if (! $found) {
        throw "Error"
    }
    $matches[1]
}

$version = Get-LatestVersion
"Latest version: $version"

if ((Test-Path $filePath) -and (Test-Path $versionFilePath)) {
    $diskVersion = Get-Content $versionFilePath
    if ($version -eq $diskVersion) {
        "Already up-to-date"
    }
}

$downloadUri = "https://github.com/labaneilers/drake/releases/download/$($version)/windows-amd64.drk.exe"
"Downloading from $downloadUri..."
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -uri $downloadUri -outfile $filePath

$version > $versionFilePath

"Ensuring $installDir is on the path..."
$path = [Environment]::GetEnvironmentVariable('PATH', [System.EnvironmentVariableTarget]::Machine);
if ($path.ToLower().Contains($installDir.ToLower()) -eq $false) {
  $path = $path + $installDir
  [System.Environment]::SetEnvironmentVariable('PATH', $path)
}