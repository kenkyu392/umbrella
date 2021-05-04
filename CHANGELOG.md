# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.0](../../releases/tag/v0.9.0) - 2021-05-04

### Added

- Added Debug provides middleware that executes the handler only if d is true.

### Changed

- Added UptimeDurationNanoseconds/Milliseconds field to Metrics.


## [0.8.0](../../releases/tag/v0.8.0) - 2021-05-03

### Added

- Added RateLimit provides middleware that limits the number of requests processed per second.
- Added RateLimitPerIP provides middleware that limits the number of requests processed per second per IP.


## [0.7.0](../../releases/tag/v0.7.0) - 2021-04-16

### Added

- Added WithRequestMetricsHookFunc option to MetricsRecorder.


## [0.6.0](../../releases/tag/v0.6.0) - 2021-04-13

### Added

- Added MetricsRecorder provides simple metrics such as request/response size and request duration.


## [0.5.0](../../releases/tag/v0.5.0) - 2021-04-04

### Added

- Added Stampede provides a simple cache middleware that is valid for a specified amount of time.


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
