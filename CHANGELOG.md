# Changelog

## [0.1.0] - 2026-03-10

### Added
- Initial implementation of BananaMCP — a single `bananamcp` binary with two subcommands:
  - `bananamcp mcp` — starts the MCP stdio server (for use with AI coding assistants)
  - `bananamcp test <prompt>` — generates an image and saves to `output.<ext>`
- `generate_image` MCP tool: generates images from a descriptive prompt and saves to a local file path
- Image generation exclusively via OpenRouter using `github.com/nealhardesty/easy-llm-wrapper`
- Default image model: `google/gemini-3.1-flash-image-preview`
- `MODEL` environment variable to override the default model
- `OPENROUTER_API_KEY` required; fail-fast if not set
- Automatic file extension derivation from response MIME type (`.png`, `.jpg`, `.webp`, `.gif`, etc.)
- `internal/generator` package with shared image generation logic
- `version.go` with semantic version `0.1.0`
- `Makefile` with `build`, `test`, `lint`, `fmt`, `tidy`, `clean`, `run-mcp`, `version`, `version-increment`, `release`, `push`, and `help` targets
- Install via `go install github.com/nealhardesty/bananamcp@latest`
