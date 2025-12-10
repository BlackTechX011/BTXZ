# Contributing Guidelines

Thank you for your interest in contributing to BTXZ. This document outlines the standards and procedures for submitting code, reporting issues, and proposing enhancements.

## Code of Conduct

All contributors are expected to maintain a professional and respectful demeanor. Harassment or abusive behavior will not be tolerated.

## Issue Reporting

### Bug Reports

When submitting a bug report, please ensure the following information is included to facilitate rapid diagnosis:

1.  **System Information**: Operating System, Architecture (amd64/arm64), and BTXZ version.
2.  **Reproduction Steps**: A concise, step-by-step guide to reproducing the issue.
3.  **Expected Behavior**: What you anticipated would happen.
4.  **Actual Behavior**: What actually happened, including full error logs and stack traces if available.

### Feature Requests

Feature proposals should include:

1.  **Use Case**: A clear explanation of the problem this feature solves.
2.  **Proposed Implementation**: A high-level technical overview of how the feature could be implemented.
3.  **Alternatives**: Other solutions that were considered.

## Development Workflow

### Pull Requests

1.  **Fork and Clone**: Fork the repository and clone it locally.
2.  **Branching**: Create a feature branch (`feature/description`) or a fix branch (`fix/issue-description`).
3.  **Coding Standards**:
    *   Adhere to standard Go idioms and formatting.
    *   Run `go fmt ./...` before committing.
    *   Ensure all new code is covered by tests where applicable.
4.  **Commit Messages**: Use imperative mood (e.g., "Add support for X", not "Added support for X").
5.  **Submission**: Open a Pull Request against the `main` branch. Provide a detailed description of the changes.

### Licensing

By submitting a contribution, you agree that your code will be licensed under the terms of the project's [LICENSE](LICENSE.md). You certify that you have the right to submit the code and that it does not infringe upon any third-party intellectual property rights.

## Security Vulnerabilities

Do not report security vulnerabilities via public GitHub issues. Please disclose them responsibly by contacting the maintainers directly.
