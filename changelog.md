# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.3+1]

### Added

- Added a changelog.
- Network state handling.
- Hook for network state.

### Changed

- Logger can now be configured via config, but also has a default logger.
- Refactored the consumer waiting method.
- Refactored logger usage to use the logger from the engine, not the global logger.

## [1.0.3] - 2025-01-24

### Changed

- Display better logging.

## [1.0.2] - 2025-01-20

### Changed

- Changing package name to "kurier"

## [1.0.1] - 2025-01-20

### Added

- Customize consumption function for flexible message processing.
- Workers pool optimization for better performance.
- Prometheus metrics for worker monitoring.

### Changed

- Refactored message consumption to use worker pool pattern.
- Improved throughput with parallel message processing.

## [1.0.0] - 2025-01-16

Handle RabbitMQ delayed message and task using Go.

We're excited to announce the first stable release of Go Delayed Message, a powerful library for handling delayed messages and tasks using RabbitMQ in Go applications.

### Requirements

- RabbitMQ 4.0.x+
- Erlang 27.x+
- Go 1.16+

### Tested Specification

- RabbitMQ 4.0.5
- Erlang 27.2
- Go 1.23.4
