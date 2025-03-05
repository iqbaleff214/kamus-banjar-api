# Changelog

## [Unreleased]

## [1.1.0] - 2025-03-06

### Added
- Added Dockerfile for containerization
- Added `.vscode` configuration for development setup
- Added endpoint for word search using Levenshtein algorithm
- Improved test cases for word search service

### Changed
- Optimized cache structure from slice/array to hash table for better search performance

### Fixed
- Fixed capitalization issue in vocabulary data (converted uppercase words to lowercase)
