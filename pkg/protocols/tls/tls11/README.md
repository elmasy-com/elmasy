# tls11

Partially implemented TLS 1.1.

Implement a part of the handshake to get information about the server.

Relevant RFC: [RFC 4346](https://datatracker.ietf.org/doc/rfc4346/)

## Errors silenced

### `readServerResponse()`

- Timeout while reading response.
- Some server send RST (TCP Reset) inmediately after Alert(Handshake Failure) at the initial handshake. This causing a `connection reset by peer` error. But if the response is exactly 7 byte, the handshake was OK, but TLS11 is not supported.

### `getSupportedCiphers()`

- Some server respond an RST without Alert(Handshake Failure) when checking only one cipher. This means that the cipher is not supported.