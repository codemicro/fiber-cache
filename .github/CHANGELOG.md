# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

(Dates are in YYYY-MM-DD format. This message is mainly for my own sake.)

## [Unreleased]

## [2.0.1] - 2020-09-14
### Fixed
* Errors are now properly forwarded to the error handler and are no longer lost in the endless void of space and time.
  * Thanks [@Fenny](https://github.com/Fenny)

## [2.0.0] - 2020-09-14
### Changed
* Update to support Fiber version 2.0.0.
### Removed
* Removed Go 1.11.x support.

## [1.0.2] - 2020-09-06
### Added
* Add support for caching without expiration.

## [1.0.1] - 2020-08-28
### Fixed
* Add missing mutex to status code storage map (t'was an oversight).

## [1.0.0] - 2020-08-27
* Initial release

[Unreleased]: https://github.com/codemicro/fiber-cache/compare/v2.0.1...HEAD
[2.0.1]: https://github.com/codemicro/fiber-cache/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/codemicro/fiber-cache/compare/v1.0.2...v2.0.0
[1.0.2]: https://github.com/codemicro/fiber-cache/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/codemicro/fiber-cache/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/codemicro/fiber-cache/releases/tag/v1.0.0
