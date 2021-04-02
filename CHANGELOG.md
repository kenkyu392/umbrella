# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0](../../releases/tag/v0.4.0) - 2021-04-03

### Added

- Added AllowAccept/DisallowAccept middleware controls the request based on the Accept header of the request.
- Added definition of the User-Agent list for possible suspicious access.


## [0.3.0](../../releases/tag/v0.3.0) - 2020-08-01

### Added

- Added CacheControl/NoCache middleware for adding Cache-Control header to response headers.
- Added AllowMethod/DisallowMethod middleware that uses request methods for access control.
- Added RealIP middleware to override RemoteAddr using X-Forwarded-For or X-Real-IP headers.
- Added Use middleware that integrates multiple middleware.

## [0.2.0](../../releases/tag/v0.2.0) - 2020-07-28

### Added

- Added AllowContentType/DisallowContentType middleware controls the request based on the Content-Type header of the request.
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
