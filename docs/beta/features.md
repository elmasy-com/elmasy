# Features of Beta 

Beta version is where the software will be usable by most people.

Below is what Elmasy should know in the beta version.

## SSL/TLS

- Versions to check: SSL3.0, TLS1.0, TLS1.1, TLS1.2 and TLS1.3
- Check cipher suites
- Check certificate
    - Validity
    - Issuer
    - Key algorithm
    - Signature algorithm
- Check TLS extensions
- StartTLS support

## FTP

- Check the used software and its version
- Check anonymous login
- Check SSL support
- FTP bounce attack

## SSH

- Check the used software and its version
- Key exchange algorithms
- Server host hey algorithms
- Encryption algorithms
- Message authentication code algorithms
- Compression algorithms

## Telnet

- Check the existence of Telnet.

## SMTP

- Check StartTLS support
- Check SSL support
- Check for open relay

## DNS

- Types:
    - A
    - AAAA
    - MX
    - TXT
    - PTR
    - CAA
    - ANY
    - AXFR

## HTTP(S)

- Check the used software and its version
- Check for forcing HTTPS
- Security headers:
    - Cache Control
    - Content-Security-Policy (CSP)
    - Strict-Transport-Security (HSTS)
    - Public Key Pinning (HPKP)
    - X-XSS-Protection
    - X-Frame-Options
    - X-Content-Type-Options
    - Feature-Policy
    - Referrer-Policy
    - Permissions-Policy

## IMAP

- Check TLS support

## POP3

- Check TLS support

## SPF

## DKIM

- Probe for popular keys (eg.: `google._domainkey`)

## DMARC

## Port scanning

- TCP: `syn` and `connect()` techniques
- UDP: probing for known protocols 
- Banner grabbing

## Subdomain enumeration

- Integrate [Amass](https://github.com/OWASP/Amass) or [Subfinder](https://github.com/projectdiscovery/subfinder)
    - **Need research**

## Database

- Store the results for statistics