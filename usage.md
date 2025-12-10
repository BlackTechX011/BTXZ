# BTXZ Usage Guide

**BTXZ** (BlackTechX Archive) is a professional-grade command-line tool designed for secure file archiving. This guide provides a comprehensive reference for all available commands, flags, and operational modes.

---

## Global Flags

These flags can be used with any command.

| Flag | Description |
| :--- | :--- |
| `--help`, `-h` | Display help information for the current command. |
| `--version`, `-v` | Show the currently installed version. |
| `--no-style` | Disable ANSI colors and rich styling (useful for scripts/logging). |

---

## Commands

### 1. `create`

Packages one or more files and/or directories into an encrypted `.btxz` archive.

**Syntax:**
```bash
btxz create [INPUTS...] -o [OUTPUT_FILE] [FLAGS]
```

**Flags:**

| Flag | Alias | Description | Required | Default |
| :--- | :--- | :--- | :--- | :--- |
| `--output` | `-o` | The destination path for the archive. | **Yes** | N/A |
| `--password` | `-p` | The encryption password. If omitted, you will be prompted securely. | No | Interactive |
| `--level` | `-l` | The hardware profile to use. Options: `low`, `default`, `max`. | No | `default` |

**Profiles:**

*   **`low` (Fast)**: Uses minimal RAM (64MB) and 1 Argon2 pass. Best for comprehensive backups on low-power devices.
*   **`default` (Balanced)**: Uses 128MB RAM. Good balance of speed and compression.
*   **`max` (Best)**: Uses 512MB RAM and 4 Argon2 passes. Maximum security against brute-force attacks and maximum compression.

**Examples:**

```bash
# Basic creation (prompts for password)
btxz create ./my-folder -o backup.btxz

# "Pro" mode: Max security, password via flag
btxz create ./database.sql -o db_backup.btxz -p "ComplexPass123!" --level max

# Archive multiple files
btxz create file1.txt file2.jpg ./folder -o mixed.btxz
```

---

### 2. `extract`

Restores files from a `.btxz` archive to the disk.

**Syntax:**
```bash
btxz extract [ARCHIVE_FILE] [FLAGS]
```

**Flags:**

| Flag | Alias | Description | Required | Default |
| :--- | :--- | :--- | :--- | :--- |
| `--output-dir` | `-o` | The directory where files will be extracted. | No | `.` (Current Dir) |
| `--password` | `-p` | The decryption password. | No | Interactive |

**Behavior:**
*   The command automatically detects whether the archive is V1, V2, or V3.
*   It performs an integrity check (MAC validation) before writing files.
*   If a file path in the archive is deemed "unsafe" (e.g., `../../etc/passwd`), it will be skipped to protect your system.

**Examples:**

```bash
# Extract to current folder
btxz extract backup.btxz

# Extract to a specific directory
btxz extract backup.btxz -o /home/user/restored
```

---

### 3. `list`

Displays the contents of an archive without extracting them. This is useful for verification or finding a specific file.

**Syntax:**
```bash
btxz list [ARCHIVE_FILE] [FLAGS]
```

**Flags:**

| Flag | Alias | Description | Required | Default |
| :--- | :--- | :--- | :--- | :--- |
| `--password` | `-p` | The decryption password. | No | Interactive |

**Note:** You must provide the correct password to list files because BTXZ encrypts the filenames and directory structure.

**Example:**
```bash
btxz list secret_files.btxz
```

---

### 4. `test`

Verifies the cryptographic integrity and compression validity of an archive.

**Syntax:**
```bash
btxz test [ARCHIVE_FILE] [FLAGS]
```

**Flags:**

| Flag | Alias | Description | Required | Default |
| :--- | :--- | :--- | :--- | :--- |
| `--password` | `-p` | The decryption password. | No | Interactive |

**What it checks:**
1.  **Authentication Tag**: Verifies that the ciphertext has not been tampered with (bit-rot or malicious editing).
2.  **Compression Stream**: Decodes the XZ stream in memory to ensure it is not corrupt.
3.  **Header Integrity**: Checks version bits and salt.

**Example:**
```bash
# Periodic backup verification script
if btxz test backup.btxz -p "pass"; then
  echo "Backup Verified"
else
  echo "Backup Corrupt!"
  exit 1
fi
```

---

### 5. `update`

Checks the official GitHub repository for a newer release and updates the `btxz` binary in-place.

**Syntax:**
```bash
btxz update
```

**Process:**
1.  Fetches `version.json` from the repository.
2.  Compares the remote version with the local version.
3.  If newer, downloads the binary for your specific OS/Arch.
4.  Replaces the current executable safely.

---

## Exit Codes

BTXZ uses standard exit codes for integration with other scripts.

*   `0`: Success. The operation completed without error.
*   `1`: General Error. (e.g., Wrong password, file not found, IO error).

