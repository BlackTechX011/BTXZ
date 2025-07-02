
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

<h1 align="center">BTXZ™ Archiver</h1>

<p align="center">
  A modern, secure, high-compression command-line archiver.
</p>

<p align="center">
    <a href="https://github.com/BlackTechX011/BTXZ/releases/latest"><img src="https://img.shields.io/github/v/release/BlackTechX011/BTXZ?style=for-the-badge&logo=github&color=blue" alt="Latest Release"></a>
    <a href="https://github.com/BlackTechX011/BTXZ/blob/main/LICENSE.md"><img src="https://img.shields.io/github/license/BlackTechX011/BTXZ?style=for-the-badge&color=lightgrey" alt="License"></a>
    <a href="https://github.com/BlackTechX011/BTXZ/actions/workflows/release.yml"><img src="https://img.shields.io/github/actions/workflow/status/BlackTechX011/BTXZ/release.yml?style=for-the-badge&logo=githubactions&logoColor=white" alt="Build Status"></a>
   
</p>

<p align="center">
  <a href="#-why-btxz">Why BTXZ?</a> •
  <a href="#-installation">Installation</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-command-reference">Commands</a> •
  <a href="#-contributing">Contributing</a> •
  <a href="#-license">License</a>
</p>

---

**BTXZ™** is a professional command-line tool for creating and extracting securely encrypted, highly compressed archives. It's built from the ground up to prioritize security and provide a polished user experience, ensuring your data is safe and the tool is a pleasure to use.

> [!NOTE]
> BTXZ encrypts **everything**, including file names and directory structures. An attacker without the password cannot learn anything about the contents of your archive.

## ✨ Why BTXZ?

| Feature | Description |
| :--- | :--- |
| **🛡️ Military-Grade Encryption** | Utilizes **AES-256-GCM** with a key derived via **Argon2id** (the winner of the Password Hashing Competition). This provides authenticated encryption, protecting against tampering and ensuring data integrity. |
| **🗜️ High-Ratio Compression** | Employs the robust **XZ (LZMA2)** algorithm to achieve superior compression ratios, saving you valuable disk space compared to standard Zip or Gzip. |
|
| **🔄 Seamless Self-Updating** | The `btxz update` command fetches the latest secure release directly from GitHub and seamlessly replaces the current executable, keeping you up-to-date with one command. |
| **🔒 Secure by Design** | Built to be resilient against malformed archives. It's hardened against path traversal attacks and will safely skip corrupted files during extraction instead of halting or crashing. |

## 🚀 Installation

The recommended way to install BTXZ™ is with our one-line installer. It automatically detects your OS and architecture, downloads the correct binary, and adds it to your system's PATH.

> [!IMPORTANT]
> The scripts below are the *only* official installation methods. Always download from the official repository to ensure you are getting a secure and untampered version of the tool.

---

### Linux / macOS / Termux

This command works on most Unix-like systems, including Debian/Ubuntu, Fedora, Arch, macOS (Intel & Apple Silicon), and Termux on Android.

```sh
curl -fsSL https://raw.githubusercontent.com/BlackTechX011/BTXZ/main/scripts/install.sh | sh
```

> [!TIP]
> After installation, you may need to restart your terminal or run `source ~/.zshrc`, `source ~/.bashrc`, etc., to refresh your `PATH` environment variable.

---

### Windows (PowerShell)

> [!NOTE]
> This command temporarily adjusts the execution policy **only for the current process**. It's a safe and standard way to run trusted remote scripts and does not permanently change your system's security settings.

**Open a new PowerShell (as a regular user) and run:**

```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass; iwr https://raw.githubusercontent.com/BlackTechX011/BTXZ/main/scripts/install.ps1 | iex
```
> [!WARNING]
> You **must** open a new PowerShell window after the installation completes. The `PATH` environment variable is only loaded when a new terminal session starts.

<details>
  <summary>Manual Installation</summary>
  
  1. Go to the [**Releases page**](https://github.com/BlackTechX011/BTXZ/releases/latest).
  2. Download the appropriate binary for your operating system and architecture (e.g., `btxz-x86_64-unknown-linux-gnu`).
  3. Rename the binary to `btxz` (or `btxz.exe` on Windows).
  4. Move the binary to a directory included in your system's `PATH` (e.g., `/usr/local/bin` on Linux/macOS, or a custom folder on Windows that you add to the Path Environment Variable).
  5. On Linux/macOS, make the binary executable: `chmod +x /usr/local/bin/btxz`.
</details>

## ⚡ Quick Start

Using BTXZ is designed to be intuitive. Here are the most common operations.

### 1. Create an Archive

To compress and encrypt files or folders into a `.btxz` archive:

```sh
# Archive a single file and a whole directory
btxz create report.docx project_assets/ -o my_archive.btxz
```

You will be prompted to enter and confirm a password securely.

> [!TIP]
> You can provide a password directly with the `-p` or `--password` flag (e.g., `btxz create ... -p "MySecret"`). This is useful for scripting but can expose the password in your shell history. Use with caution.

### 2. Extract an Archive

To decompress and decrypt an archive:

```sh
# Extract to the current directory
btxz extract my_archive.btxz

# Extract to a specific output directory
btxz extract my_archive.btxz -o ./restored_files
```

You will be prompted for the password.

### 3. List Archive Contents

To see the file and directory structure inside an archive without extracting it:

```sh
btxz list my_archive.btxz
```

This requires the password, as the file list itself is encrypted.

## 📖 Command Reference

<details>
  <summary>Click to expand the full command reference</summary>

| Command | Alias | Description | Options |
| :--- | :--- | :--- | :--- |
| `create` | `c` | Creates a new encrypted, compressed archive. | `-o, --output <path>` (Required)<br>`-p, --password <pass>` |
| `extract`| `x` | Extracts files and folders from an archive. | `-o, --output <path>`<br>`-p, --password <pass>` |
| `list`   | `l` | Lists the contents of an archive. | `-p, --password <pass>` |
| `update` | `u` | Checks for and installs the latest version of BTXZ. | |
| `help`   | | Displays help information for a command. | |

</details>

---

## 🗺️ Project Roadmap

This project is actively developed. Here is a list of planned features. Contributions are welcome!

### Core Features
- [x] Core `create`, `extract`, and `list` commands
- [x] Secure password-based encryption (AES-256-GCM + Argon2id)
- [x] High-ratio compression with XZ
- [ ] Add support for different compression levels (`--level`)
- [ ] Implement a `modify` command to add/remove files from existing archives
- [ ] Implement a `test` command to verify archive integrity

### User Experience
- [x] Automated release workflow with cross-compiled binaries
- [x] Self-update mechanism (`btxz update`)
- [ ] Add a global configuration file (`~/.config/btxz/config.json`)
- [ ] Implement interactive mode (`btxz --interactive`) for guided operations
- [ ] Add detailed progress bars for large file operations
- [ ] Implement a GUI wrapper for the CLI tool

### Documentation & Community
- [x] `LICENSE.md` with custom EULA
- [x] `CONTRIBUTING.md` and Issue Templates
- [ ] Add advanced usage examples and a FAQ to this README
- [ ] Create a GitHub Pages site for full documentation

> **Have an idea or found a bug?** [**Open an issue!**](https://github.com/BlackTechX011/BTXZ/issues/new/choose) We'd love to hear from you.

## 🤝 Contributing

Contributions are the backbone of open source. We welcome contributions of all kinds, from filing detailed bug reports to implementing new features.

Before you start, please take a moment to read our guidelines:

-   **[Contribution Guide](CONTRIBUTING.md):** The main guide for how to submit pull requests, our coding standards, and the development process.
-   **[Open an Issue](https://github.com/BlackTechX011/BTXZ/issues/new/choose):** The best place to report a bug, ask a question, or propose a new feature.

## 🛡️ Security Model

> [!CAUTION]
> This software is provided "as is" without warranty of any kind. While it is designed with strong security principles, you are responsible for securely managing your passwords. **There is no way to recover a lost password.**

The security of BTXZ™ is a top priority. If you discover a security vulnerability, we ask that you report it to us privately to protect our users.

**Please do not open a public GitHub issue for security-related concerns.**

Instead, send a detailed report directly to:
**`BlackTechX@proton.me`**

We will make every effort to respond to your report in a timely manner.

## ⚖️ License

This software is distributed under a custom End-User License Agreement (EULA).

> [!IMPORTANT]
> The license grants permission for **personal, non-commercial use only**. For any other use, including commercial, corporate, or government, please contact the author.

Please see the [**LICENSE.md**](LICENSE.md) file for the full terms and conditions.

---
*BTXZ™ is a trademark of [BlackTechX011](https://github.com/BlackTechX011). All rights reserved.*

<p align="right">(<a href="#top">back to top</a>)</p>