# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| main    | ✅ |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in Poll Creator, please follow these steps:

### 🔒 Private Disclosure

1. **DO NOT** create a public GitHub issue
2. Email security details to: [your-email@example.com] (replace with your actual email)
3. Include as much detail as possible:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if you have one)

### ⏱️ Response Timeline

- **Initial Response**: Within 48 hours
- **Vulnerability Assessment**: Within 1 week
- **Fix Timeline**: Depends on severity
  - Critical: 1-3 days
  - High: 1 week
  - Medium: 2 weeks
  - Low: 1 month

### 🛡️ Security Measures

This project implements several security measures:

- Input validation on all API endpoints
- CORS configuration
- Rate limiting considerations
- Container security best practices
- Dependency vulnerability scanning via GitHub Dependabot

### 📋 Security Checklist for Contributors

When contributing code, please ensure:

- [ ] User inputs are properly validated and sanitized
- [ ] No sensitive information is logged or exposed
- [ ] Authentication/authorization is properly implemented
- [ ] Dependencies are up to date and secure
- [ ] Error messages don't leak sensitive information

### 🔍 Automated Security

- **GitHub Dependabot**: Automatically checks for vulnerable dependencies
- **CodeQL Analysis**: Static code analysis for security issues
- **Docker Image Scanning**: Container vulnerability scanning

## Acknowledgements

We appreciate security researchers who responsibly disclose vulnerabilities. Contributors will be acknowledged in our security acknowledgments section (unless they prefer to remain anonymous).

---

Thank you for helping keep Poll Creator secure! 🔒