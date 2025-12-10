<#
.SYNOPSIS
    Professional Installer for BTXZ Archiver (Windows)
.DESCRIPTION
    Downloads, verifies, and installs the latest version of BTXZ.
    Adds the installation directory to the User PATH if needed.
#>

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "BlackTechX011/BTXZ"
$VersionUrl = "https://raw.githubusercontent.com/$Repo/main/version.json"
$InstallDir = "$env:LOCALAPPDATA\Programs\btxz"
$ExeName = "btxz.exe"

# Colors
function Print-Header($Msg) { Write-Host -ForegroundColor Cyan "`n$Msg" }
function Print-Success($Msg) { Write-Host -ForegroundColor Green "SUCCESS: $Msg" }
function Print-Info($Msg) { Write-Host -ForegroundColor Gray "INFO: $Msg" }
function Print-Error($Msg) { Write-Host -ForegroundColor Red "ERROR: $Msg"; exit 1 }

# --- BANNER ---
Write-Host -ForegroundColor Cyan @"
██████╗ ████████╗██╗  ██╗███████╗
██╔══██╗╚══██╔══╝╚██╗██╔╝╚══███╔╝
██████╔╝   ██║    ╚███╔╝   ███╔╝ 
██╔══██╗   ██║    ██╔██╗  ███╔╝  
██████╔╝   ██║   ██╔╝ ██╗███████╗
╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝
PROFESSIONAL SECURE ARCHIVER
"@

# 1. Fetch Version Info
Print-Header "CHECKING UPDATES"
try {
    $JsonContent = Invoke-RestMethod -Uri $VersionUrl
} catch {
    Print-Error "Failed to fetch version info: $_"
}

$LatestVersion = $JsonContent.version
$PlatformKey = "windows-amd64" # BTXZ only targets 64-bit windows for now
$PlatformInfo = $JsonContent.platforms.$PlatformKey

if (-not $PlatformInfo) {
    Print-Error "No Windows release found in manifest."
}

$DownloadUrl = $PlatformInfo.url
$ExpectedHash = $PlatformInfo.sha256

Print-Info "Found Version: v$LatestVersion"

# 2. Prepare Install Directory
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

# 3. Download
Print-Header "DOWNLOADING"
$OutFile = Join-Path $InstallDir $ExeName
Print-Info "Source: $DownloadUrl"
Print-Info "Dest:   $OutFile"

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $OutFile
} catch {
    Print-Error "Download failed: $_"
}

# 4. Verify Checksum
if ($ExpectedHash) {
    Print-Header "VERIFYING SECURITY"
    $FileHash = (Get-FileHash $OutFile -Algorithm SHA256).Hash.ToLower()
    
    if ($FileHash -eq $ExpectedHash.ToLower()) {
        Print-Success "SHA256 Checksum Verified."
    } else {
        Remove-Item $OutFile -ErrorAction SilentlyContinue
        Write-Host "Expected: $ExpectedHash"
        Write-Host "Actual:   $FileHash"
        Print-Error "Checksum mismatch! The file may be corrupted or tampered with."
    }
} else {
    Write-Host -ForegroundColor Yellow "WARNING: No checksum provided in manifest."
}

# 5. Add to PATH
Print-Header "CONFIGURING ENVIRONMENT"
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Print-Info "Adding $InstallDir to User PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $Env:Path += ";$InstallDir"
    Print-Success "Path updated."
} else {
    Print-Info "Path already configured."
}

# 6. Final Check
Print-Header "VERIFICATION"
try {
    $InstalledVersion = & $OutFile --version
    Print-Success "Installed $InstalledVersion successfully!"
    Write-Host "`nYou can now use 'btxz' from any new terminal window."
} catch {
    Print-Error "Binary check failed. The file might not be executable."
}
