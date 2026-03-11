# Changelog

## [0.2.0] - 2026-03-11

### Added
- Google Vertex AI backend for image generation, activated via `--vertex` flag
- `google.golang.org/genai` SDK integration for Vertex AI image generation
- `ImageGenerator` interface to abstract backend selection (OpenRouter vs Vertex AI)
- Factory function `generator.New(useVertex)` for backend instantiation
- `GCLOUD_PROJECT_ID` and `GCLOUD_LOCATION` environment variables for Vertex AI configuration
- Vertex AI MCP server configuration example in README

### Changed
- Refactored `internal/generator` package: split monolithic `generator.go` into interface (`generator.go`), OpenRouter backend (`openrouter.go`), and Vertex AI backend (`vertex.go`)
- `main.go` now parses a `--vertex` global flag and passes an `ImageGenerator` to `runMCP`/`runTest`
- Default model is `gemini-2.5-flash-image` (Nano Banana) for Vertex AI and `google/gemini-3.1-flash-image-preview` (Nano Banana 2) for OpenRouter
- `runMCP` and `runTest` now accept an `ImageGenerator` parameter instead of calling the package-level function directly
- Moved `OPENROUTER_API_KEY` validation from `main.go` into the OpenRouter generator constructor
- Go minimum version bumped to 1.24 (required by genai SDK)

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
