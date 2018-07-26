
$installDir = Join-Path $env:ALLUSERSPROFILE "drake"
if (![System.IO.Directory]::Exists($installDir)) {[void][System.IO.Directory]::CreateDirectory($installDir)}
$filePath = Join-Path $installDir "drk.exe"

# $r = Invoke-RestMethod -uri "https://api.github.com/repos/labaneilers/drake/releases/latest"
# $downloadUri = $r.assets | ? { $_.name -match 'windows-amd64' } | % { $_.browser_download_url }

function Get-LatestVersion() {
    $r = (iwr "https://github.com/labaneilers/drake/releases").Content
    $line = $r.Split("`n") | ? { $_ -match "download.*windows-amd64" }
    $found = $line -match 'download\/(.*)\/windows-amd64'
    if (! $found) {
        throw "Error"
    }
    
    $matches[1]
}

$version = Get-LatestVersion
$downloadUri = "https://github.com/labaneilers/drake/releases/download/$($version)/windows-amd64.drk.exe"


[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -uri $downloadUri -outfile $filePath

Write-Output 'Ensuring drk is on the path'
$path = [Environment]::GetEnvironmentVariable('PATH', [System.EnvironmentVariableTarget]::Machine);
if ($path.ToLower().Contains($installDir.ToLower()) -eq $false) {
  $path = $path + $installDir
  [System.Environment]::SetEnvironmentVariable('PATH', $path)
}