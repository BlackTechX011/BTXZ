<#
.SYNOPSIS
    Installs the BTXZ™ command-line tool for Windows.
.DESCRIPTION
    This script automatically detects the Windows OS version to download the
    latest "modern" (Windows 11+) or "compat" (Windows 10 and older)
    binary for BTXZ from GitHub. It installs it to the user's home
    directory (~/.btxz) and adds the installation directory to the user's PATH.
.EXAMPLE
    To run from the web (ensure Execution Policy is set appropriately):
    iwr https://raw.githubusercontent.com/BlackTechX011/BTXZ/main/scripts/install.ps1 | iex
.NOTES
    Author: BlackTechX011
    License: BTXZ EULA (https://github.com/BlackTechX011/BTXZ/blob/main/LICENSE.md)
    Modified to include OS version detection for modern/compat builds.
#>

# Stop script on any error
$ErrorActionPreference = 'Stop'

# --- Configuration ---
$Repo = "BlackTechX011/BTXZ"
$InstallDir = Join-Path $HOME ".btxz"
$ExeName = "btxz.exe"

# --- Main Logic ---
try {
    Write-Host "Starting BTXZ™ installation for Windows..." -ForegroundColor Cyan

    # 1. Detect Architecture
    if ($env:PROCESSOR_ARCHITECTURE -ne "AMD64") {
        throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE). BTXZ currently supports AMD64 (64-bit) on Windows."
    }
    $Arch = "amd64"
    $Os = "windows"

    # --- NEW: OS Version Detection to select correct binary ---
    # Windows 11 has a build number of 22000 or greater.
    $OsVersion = [System.Environment]::OSVersion.Version
    $BuildNumber = $OsVersion.Build

    if ($BuildNumber -ge 22000) {
        Write-Host "Detected Windows 11 or newer (Build $BuildNumber)." -ForegroundColor Green
        $BinaryName = "btxz-$Os-$Arch-modern.exe"
    } else {
        Write-Host "Detected Windows 10 or older (Build $BuildNumber)." -ForegroundColor Green
        $BinaryName = "btxz-$Os-$Arch-compat.exe"
    }
    # --- END NEW ---

    Write-Host "Detected System: $Os-$Arch. Target binary: $BinaryName" -ForegroundColor Cyan

    # 2. Get the download URL for the latest release
    $ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"
    Write-Host "Fetching latest release information from GitHub..." -ForegroundColor Cyan

    $latestRelease = Invoke-RestMethod -Uri $ApiUrl
    $asset = $latestRelease.assets | Where-Object { $_.name -eq $BinaryName }

    if (-not $asset) {
        throw "Could not find a download URL for '$BinaryName'. The release may be missing."
    }

    $DownloadUrl = $asset.browser_download_url
    Write-Host "Download URL: $DownloadUrl" -ForegroundColor Cyan

    # 3. Download and Install
    $InstallPath = Join-Path $InstallDir $ExeName
    Write-Host "Installing to $InstallPath..." -ForegroundColor Cyan

    # Ensure the installation directory exists
    if (-not (Test-Path -Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir | Out-Null
    }

    Invoke-WebRequest -Uri $DownloadUrl -OutFile $InstallPath

    # 4. Add to PATH (permanently for the current user)
    Write-Host "Adding installation directory to your PATH..." -ForegroundColor Cyan
    $UserPath = [System.Environment]::GetEnvironmentVariable('Path', 'User')

    if (-not ($UserPath -split ';' -contains $InstallDir)) {
        # Using an array to avoid potential double semicolons
        $PathArray = $UserPath -split ';' | Where-Object { $_ -ne '' }
        $PathArray += $InstallDir
        $NewPath = $PathArray -join ';'
        
        [System.Environment]::SetEnvironmentVariable('Path', $NewPath, 'User')

        # Also update the current session's PATH
        $env:PATH += ";$InstallDir"

        Write-Host "Updated user PATH. Please restart your terminal for the changes to take full effect." -ForegroundColor Yellow
    } else {
        Write-Host "PATH is already configured." -ForegroundColor Green
    }

    Write-Host ""
    Write-Host "BTXZ™ was installed successfully!" -ForegroundColor Green
    Write-Host "You can now run 'btxz' from a new terminal window." -ForegroundColor Green

}
catch {
    Write-Host "An error occurred during installation:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
