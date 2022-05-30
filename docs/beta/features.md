# Features of Beta 

Beta version is between `v0.1` and `v1.0`.

Below is what Elmasy should know until v1.0.

## SSL/TLS

- [x] Versions to check:
    - [x] SSL3.0
    - [x] TLS1.0
    - [x] TLS1.1
    - [x] TLS1.2
    - [x] TLS1.3
- [x] Check cipher suites
- [x] Check certificate
- [ ] Check TLS extensions
- [ ] StartTLS support

## FTP

- [ ] Check the used software and its version
- [ ] Check anonymous login
- [ ] Check SSL support
- [ ] FTP bounce attack

## SSH

- [ ] Check the used software and its version
- [ ] Key exchange algorithms
- [ ] Server host hey algorithms
- [ ] Encryption algorithms
- [ ] Message authentication code algorithms
- [ ] Compression algorithms

## Telnet

- [ ] Check the existence of Telnet.

## SMTP

- [ ] Check StartTLS support
- [ ] Check SSL support
- [ ] Check for open relay

## DNS

- [ ] Types:
    - [x] A
    - [x] AAAA
    - [x] MX
    - [x] TXT
    - [ ] PTR
    - [ ] CAA
    - [ ] ANY
    - [ ] AXFR

## HTTP(S)

- [ ] Check the used software and its version
- [ ] Check for forcing HTTPS
- [ ] Security headers:
    - [ ] Cache Control
    - [ ] Content-Security-Policy (CSP)
    - [ ] Strict-Transport-Security (HSTS)
    - [ ] Public Key Pinning (HPKP)
    - [ ] X-XSS-Protection
    - [ ] X-Frame-Options
    - [ ] X-Content-Type-Options
    - [ ] Feature-Policy
    - [ ] Referrer-Policy
    - [ ] Permissions-Policy

## IMAP

- [ ] Check TLS support

## POP3

- [ ] Check TLS support

## SPF

- [ ] Check

## DKIM

- [ ] Probe for popular keys (eg.: `google._domainkey`)

## DMARC

- [ ] Check

## Port scanning

- [x] TCP
    - [x] `syn`
    - [x] `connect()`
- [ ] UDP: probing for known protocols 
- [ ] Banner grabbing

## Subdomain enumeration

- [ ] Integrate [Amass](https://github.com/OWASP/Amass)

## Database

- [ ] Store the results for statistics