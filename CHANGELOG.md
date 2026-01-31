# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-01-31

### Added
- Public API release - framework can now be imported by external projects
- Core interfaces package (`layout`) with Focusable, Sizeable, Positional, Help interfaces
- UI components package (`ui`) with:
  - Theme system with HCL color blending and gradients
  - Virtualized list components with filtering and selection
  - Dialog system (info, confirm, input, select list)
  - Layout utilities (flex, grid, split, area)
  - Rendering engine abstraction (Bubbletea, Ultraviolet, Direct)
  - Component library (files, header, messages, progress, sidebar, status, treefiles)
- TUI framework package (`tui`) with:
  - Application framework with page and dialog management
  - Lifecycle management
  - Syntax highlighting with Chroma
  - Animated components
  - High-level components (completions, dialogs, logo, messages, image)
  - Experimental diff viewer and list implementations

### Changed
- **BREAKING**: Moved packages from `internal/` to top-level for public API
  - `internal/layout` → `layout`
  - `internal/ui` → `ui`
  - `internal/tui` → `tui`
- Updated 74 import paths across all examples and test files
- Improved DiffViewer component with proper scrolling logic
- Enhanced test coverage with 26 passing tests

### Fixed
- Fixed scrolling tests in `tui/exp/diffview` to use correct content height
- Fixed import paths in all example programs

### Removed
- Removed non-functional `examples/tasks-demo` (referenced non-existent components)

## [0.9.0] - Previous

### Added
- Initial framework extraction from Crush CLI
- Core component interfaces
- Theme system with gradients
- Dialog system
- Virtualized lists
- File picker
- Message components
- Markdown rendering
- Syntax highlighting

---

[1.0.0]: https://github.com/wwsheng009/taproot/releases/tag/v1.0.0
[0.9.0]: https://github.com/wwsheng009/taproot/releases/tag/v0.9.0
