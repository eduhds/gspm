# Install gspm on Windows
# Command:
# powershell -c "irm https://dub.sh/gspm.ps1 | iex"
# ------------------------------------------------

Write-Host "Installing gspm..."

$appVersion = "v1.0.1"

# Setup directory
# ------------------------------------------------

$targetDir = "C:\Program Files\gspm"

if (!(Test-Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir
}

# Setup Path
# ------------------------------------------------

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$userPath >> C:\Path_Before_gspm.bkp.txt

if (!($userPath -Match $targetDir)) {
    # Current session
    $env:Path += ";$targetDir"

    # Permanent
    $userPath += ";$targetDir"
    [Environment]::SetEnvironmentVariable("Path", $userPath, "User")   
}

# Download gspm
# ------------------------------------------------

$appFile = "gspm_" + $appVersion + "_Windows_x86_64.zip"
$appUrl = "https://github.com/eduhds/gspm/releases/download/$appVersion/$appFile"
Invoke-WebRequest $appUrl -OutFile $appFile

if (Test-Path $appFile) {
    Expand-Archive $appFile -DestinationPath $targetDir
    Remove-Item $appFile
}

gspm --version

Write-Host "✅ Done."
