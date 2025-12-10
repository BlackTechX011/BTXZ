
<div id="top"></div>
<p align="center">
  <pre>
██████╗ ████████╗██╗  ██╗███████╗
██╔══██╗╚══██╔══╝╚██╗██╔╝╚══███╔╝
██████╔╝   ██║    ╚███╔╝   ███╔╝ 
██╔══██╗   ██║    ██╔██╗  ███╔╝  
██████╔╝   ██║   ██╔╝ ██╗███████╗
╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝
  </pre>
</p>

<h1 align="center">BTXZ Archiver</h1>

<p align="center">
  <strong>Secure. Adaptive. Professional.</strong>
</p>

<p align="center">
    <a href="https://github.com/BlackTechX011/BTXZ/releases/latest"><img src="https://img.shields.io/github/v/release/BlackTechX011/BTXZ?style=flat-square&color=0055ff" alt="Latest Release"></a>
    <a href="https://github.com/BlackTechX011/BTXZ/blob/main/LICENSE.md"><img src="https://img.shields.io/github/license/BlackTechX011/BTXZ?style=flat-square&color=lightgrey" alt="License"></a>
    <a href="https://github.com/BlackTechX011/BTXZ/actions/workflows/release.yml"><img src="https://img.shields.io/github/actions/workflow/status/BlackTechX011/BTXZ/release.yml?style=flat-square" alt="Build Status"></a>
</p>

---

## Overview

**BTXZ** is a high-performance command-line archiving utility designed for environments where security and data integrity are paramount. It combines state-of-the-art cryptography with industry-leading compression algorithms to force a strict "security-first" posture.

Unlike general-purpose archivers, BTXZ assumes all data is sensitive. It employs authenticated encryption for both file contents and metadata (filenames, directory structure), ensuring zero information leakage to unauthorized parties.

## Technical Specifications (V3 Protocol)

The V3 format is the current standard for all new BTXZ archives. It is engineered to withstand modern cryptographic attacks, including nonce reuse vulnerabilities and brute-force attempts.

| Component | Specification | Details |
| :--- | :--- | :--- |
| **Encryption** | **XChaCha20-Poly1305** | Authenticated Encryption with Associated Data (AEAD). Uses a **24-byte nonce** to eliminate the risk of random nonce collisions, making it safe for massive scale deployments. |
| **Key Derivation** | **Argon2id** | The winner of the Password Hashing Competition. Configured with adaptive memory hardness (up to 512MB) to render GPU-based brute-force attacks infeasible. |
| **Compression** | **LZMA2 (XZ)** | High-ratio compression algorithm optimized for binary data and text. Tuned with variable dictionary sizes based on the selected hardware profile. |
| **Container** | **Tar** | POSIX-compliant tar stream enabling preservation of permissions, ownership, and directory structures. |

## Adaptive Hardware Profiles

BTXZ V3 introduces **Adaptive Profiles**, allowing the operator to tailor the cryptographic and compression workload to the available hardware resources.

| Profile | Flag | Memory Requirement | Description |
| :--- | :--- | :--- | :--- |
| **Low** | `--level low` | ~64 MB | Optimized for embedded systems (Raspberry Pi), mobile devices in Termux, and legacy hardware. Fast encryption, moderate compression. |
| **Balanced** | `--level default` | ~128 MB | The standard profile. Balances high compression ratios with reasonable memory usage. Suitable for most desktops and servers. |
| **Max** | `--level max` | ~512 MB | **Paranoid Mode.** Maximizes encryption difficulty (4 passes of Argon2id) and uses Ultra-level XZ compression settings. Recommended for archival of critical data on powerful hardware. |

## Installation

### Automated Installer

The following command automatically detects the operating system (Windows, Linux, macOS) and architecture (amd64, arm64), downloads the latest binary, and installs it to the system path.

**Linux / macOS / Termux**
```sh
curl -fsSL https://raw.githubusercontent.com/BlackTechX011/BTXZ/main/scripts/install.sh | sh
```

**Windows (PowerShell)**
```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass; iwr https://raw.githubusercontent.com/BlackTechX011/BTXZ/main/scripts/install.ps1 | iex
```

### Manual Installation

1. Navigate to the [Releases](https://github.com/BlackTechX011/BTXZ/releases/latest) page.
2. Download the binary corresponding to your OS and CPU architecture.
3. Verify the SHA256 checksum against the `sha256sums.txt` file provided in the release.
4. Move the executable to a directory in your system `PATH`.

## Usage Guide

### Creating a Secure Archive

To create a new archive, use the `create` command. You must specify the output path via `-o`.

```sh
# Create a standard archive (Balanced Profile)
btxz create ./documents ./projects -o confidential.btxz

# Create a maximium security archive for long-term storage
btxz create ./database_dump.sql -o backup.btxz --level max
```

The system will securely prompt for a password if one is not provided via flags.

### Extracting an Archive

Extraction automatically detects the archive version (V1, V2, or V3) and applies the correct decryption routine.

```sh
# Extract to the current directory
btxz extract archive.btxz

# Extract to a specific destination
btxz extract archive.btxz -o /path/to/destination
```

### Verifying Integrity

The `test` command validates the integrity of the archive structure and the authenticity of the ciphertext without writing data to disk. This ensures the file has not been tampered with or corrupted.

```sh
btxz test archive.btxz
```

### Listing Contents

Operators can inspect the contents of an encrypted archive without extraction. This requires the valid decryption key.

```sh
btxz list archive.btxz
```

## Security Model

**Authentication is Mandatory:** BTXZ uses AEAD (Authenticated Encryption). Any modification to the ciphertext (bit-flipping, truncation) will be detected during decryption, and the operation will be aborted immediately.

**Metadata Protection:** The file manifest is encrypted. An attacker accessing the `.btxz` file cannot determine the filenames, directory structure, or file sizes contained within.

**Memory Hardness:** The Argon2id parameters are tuned to exceed the memory bandwidth available to typical ASIC/FPGA cracking rigs, forcing attackers to use general-purpose RAM, which significantly increases the cost of an attack.

## License

This software is distributed under a proprietary End-User License Agreement (EULA). It is free for personal, non-commercial use. Commercial deployment requires prior authorization. See `LICENSE.md` for full terms.

Copyright © 2025 BlackTechX011. All Rights Reserved.
