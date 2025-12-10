# Security Policy

## Supported Versions

We currently support the following versions of BTXZ. We strongly recommend always running the latest version to ensure you have the most up-to-date security patches.

| Version | Supported          |
| ------- | ------------------ |
| v1.2.x  | :white_check_mark: |
| v1.1.x  | :x:                |
| v1.0.x  | :x:                |
| < v1.0  | :x:                |

## Reporting a Vulnerability

We take the security of BTXZ seriously. If you have discovered a vulnerability, we appreciate your help in disclosing it to us in a responsible manner.

### Process

1.  **Do not open a public GitHub issue.** This allows us to patch the vulnerability before it can be exploited.
2.  **Email us directly** at [BlackTechX@proton.me](mailto:BlackTechX@proton.me).
    *   If possible, please encrypt your message using our PGP key (Key ID: `0xXYZ...`).
3.  Include a detailed description of the vulnerability, steps to reproduce it, and any proof-of-concept code.

### Our Response

1.  We will acknowledge receipt of your report within 48 hours.
2.  We will investigate the issue and confirm its validity.
3.  We will provide an estimated timeline for a fix.
4.  We will credit you (if desired) in the release notes once the patch is published.

## Security Features

BTXZ relies on well-tested cryptographic primitives:

*   **AEAD**: XChaCha20-Poly1305 (Go standard library `golang.org/x/crypto/chacha20poly1305`)
*   **KDF**: Argon2id (Go standard library `golang.org/x/crypto/argon2`)
*   **Compression**: LZMA2 (via `github.com/ulikunitz/xz`)

We perform regular audits of our dependency tree to ensure no supply chain vulnerabilities are introduced.
