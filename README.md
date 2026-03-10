# bananamcp

A lightweight [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server written in Go. It communicates via stdio and exposes a single `generate_image` tool that lets AI coding assistants (Cursor, Claude Code, etc.) generate images via [OpenRouter](https://openrouter.ai/) and save them directly to the local filesystem.

## Features

- Single MCP tool: `generate_image`
- Image generation exclusively via OpenRouter using `github.com/nealhardesty/easy-llm-wrapper`
- Default model: `google/gemini-3.1-flash-image-preview`
- Automatic file extension detection from response MIME type
- Fail-fast if `OPENROUTER_API_KEY` is not set
- `bananamcp test <prompt>` subcommand for quick manual testing

## Requirements

- Go 1.23+
- `OPENROUTER_API_KEY` environment variable

## Installation

```bash
go install github.com/nealhardesty/bananamcp@latest
```

This installs the `bananamcp` binary to `$(go env GOPATH)/bin`.

Or build from source:

```bash
make build
```

## Configuration

| Variable             | Required | Default                                   | Description                        |
|----------------------|----------|-------------------------------------------|------------------------------------|
| `OPENROUTER_API_KEY` | Yes      | —                                         | OpenRouter API key                 |
| `MODEL`              | No       | `google/gemini-3.1-flash-image-preview`   | Override the image generation model |

## Usage

```
bananamcp mcp              Start the MCP stdio server
bananamcp test <prompt>    Generate an image and save to output.<ext>
bananamcp --version        Print version and exit
```

## MCP Server

The server communicates over stdio using the Model Context Protocol. Configure it in your MCP client (e.g., Cursor, Claude Code):

```json
{
  "mcpServers": {
    "bananamcp": {
      "command": "/path/to/bananamcp",
      "args": ["mcp"],
      "env": {
        "OPENROUTER_API_KEY": "your-key-here"
      }
    }
  }
}
```

### Tool: `generate_image`

**Description:** Generates an image based on a descriptive prompt and saves it to the specified local file path.

**Parameters:**

| Name        | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `prompt`    | string | Yes      | The detailed visual description and stylistic constraints for the image. |
| `save_path` | string | Yes      | The base path where the image file should be saved. The file extension is derived from the response MIME type (`.png`, `.jpg`, `.webp`, etc.). If `save_path` already has an extension it is used as-is. |

**Returns:** `"Successfully generated and saved image to: <final_path>"` on success, or a human-readable error message on failure.

## Test Mode

Generate an image directly from the command line:

```bash
bananamcp test A maine coon cat on a fancy throne holding a beer.
# Saves to output.png (or .jpg/.webp — derived from model response)
```

## Development

```bash
make help          # Show all targets
make build         # Build bananamcp binary
make test          # Run tests with race detector
make lint          # Run go vet (and golangci-lint if installed)
make fmt           # Format code
make tidy          # go mod tidy
make clean         # Remove build artifacts
make version       # Show current version
```

## Architecture

```
bananamcp/
├── main.go                    # Single entry point (mcp + test subcommands)
├── version.go                 # Semantic version
├── internal/
│   └── generator/
│       └── generator.go       # Image generation logic
├── Makefile
└── go.mod
```

The core business logic lives in `internal/generator`. The `bananamcp` binary provides both `mcp` (stdio MCP server) and `test` (CLI image generation) subcommands.

## License

MIT
