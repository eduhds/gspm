Write-Host "Installing gspm..."

$appVersion = "0.2.3"

# Setup directory
# ------------------------------------------------

$targetDir = "C:\Program Files\gspm"

if (!(Test-Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir
}

# Setup Path
# ------------------------------------------------

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$userPath >> C:\Path.bkp.txt

# Current session
$env:Path += ";$targetDir"

# Permanent
$userPath += ";$targetDir"
[Environment]::SetEnvironmentVariable("Path", $userPath, "User")

# Download gspm
# ------------------------------------------------

$appFile = "gspm_Windows_x86_64.zip"
$appUrl = "https://github.com/eduhds/gspm/releases/download/v$appVersion/$appFile"
Invoke-WebRequest $appUrl -OutFile $appFile

if (Test-Path $appFile) {
    Expand-Archive $appFile -DestinationPath $targetDir
    Remove-Item $appFile
}

gspm --version

Write-Host "✅ Done."
