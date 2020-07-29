# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0](../../releases/tag/v0.2.0) - 2020-07-28

### Added

- Added AllowContentType/DisallowContentType middleware for authentication using the Content-Type header.
- Added Timeout middleware for cancelling the context of the request scope.
- Added RequestHeader/ResponseHeader middleware for editing request and response headers.
- Added HSTS middleware for adding Strict-Transport-Security header to response headers.
- Added Clickjacking middleware for adding X-Frame-Options header to response headers.
- Added ContentSniffing middleware for adding X-Content-Type-Options header to response headers.

## [0.1.0](../../releases/tag/v0.1.0) - 2020-07-23

### Added

- Added AllowHTTPHeader/DisallowHTTPHeader for authentication using headers.
- Added AllowUserAgent/DisallowUserAgent middleware for authentication using the User-Agent header.
- Added Context middleware for editing the context of the request scope.