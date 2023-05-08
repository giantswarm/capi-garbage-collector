# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2023-05-08

### Changed

- Use patch when deleting the finalizer from the `MachinePools`.
- Use right format on dates when logging.
- Don't requeue `MachinePools` that are not being deleted. When they are deleted, they will be queued anyway.
- Give RBAC permissions to watch `Secrets` to avoid `client-go` errors when managing the client cache for `Secrets`.

## [0.1.5] - 2023-02-01

### Added

- Allow the capi-garbage-collector SA to create events 

## [0.1.4] - 2022-11-24

### Fixed

- Increased memory limit and added VPA

## [0.1.3] - 2022-11-21

### Fixed

- Added missing RBAC permission for leases

## [0.1.2] - 2022-11-21

### Fixed

- Add experimental api to scheme

## [0.1.1] - 2022-11-21

### Fixed

- Remove unneeded args

## [0.1.0] - 2022-11-21

- Initial release

[Unreleased]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.1.5...v0.2.0
[0.1.5]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/giantswarm/capi-garbage-collector/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/capi-garbage-collector/releases/tag/v0.1.0
