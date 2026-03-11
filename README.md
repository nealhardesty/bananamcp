# bananamcp

A lightweight [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server written in Go. It communicates via stdio and exposes a single `generate_image` tool that lets AI coding assistants (Cursor, Claude Code, etc.) generate images and save them directly to the local filesystem.

Supports two image generation backends:
- **OpenRouter** (default) — via `github.com/nealhardesty/easy-llm-wrapper`
- **Google Vertex AI** — via the `google.golang.org/genai` SDK (activated with `--vertex`)

## Features

- Single MCP tool: `generate_image`
- Dual backend support: OpenRouter and Google Vertex AI
- Default model: `gemini-3.1-flash-image-preview` (Nano Banana 2)
- Automatic file extension detection from response MIME type
- Fail-fast on missing required environment variables
- `bananamcp test <prompt>` subcommand for quick manual testing

## Requirements

- Go 1.24+
- **OpenRouter mode** (default): `OPENROUTER_API_KEY` environment variable
- **Vertex AI mode** (`--vertex`): `GCLOUD_PROJECT_ID` environment variable and `gcloud auth application-default login` completed

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

### OpenRouter mode (default)

| Variable             | Required | Default                                   | Description                        |
|----------------------|----------|-------------------------------------------|------------------------------------|
| `OPENROUTER_API_KEY` | Yes      | —                                         | OpenRouter API key                 |
| `MODEL`              | No       | `google/gemini-3.1-flash-image-preview`   | Override the image generation model |

### Vertex AI mode (`--vertex`)

| Variable             | Required | Default                            | Description                        |
|----------------------|----------|------------------------------------|------------------------------------|
| `GCLOUD_PROJECT_ID`  | Yes      | —                                  | Google Cloud project ID            |
| `GCLOUD_LOCATION`    | No       | `us-central1`                      | Google Cloud region                |
| `MODEL`              | No       | `gemini-2.5-flash-image`           | Override the image generation model |

**Prerequisites for Vertex AI:** Run `gcloud auth application-default login` before using `--vertex` mode. The SDK uses Application Default Credentials for authentication.

## Usage

```
bananamcp [--vertex] mcp              Start the MCP stdio server
bananamcp [--vertex] test <prompt>    Generate an image and save to output.<ext>
bananamcp --version                   Print version and exit
```

## MCP Server

The server communicates over stdio using the Model Context Protocol. Configure it in your MCP client (e.g., Cursor, Claude Code):

### OpenRouter (default)

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

### Vertex AI

```json
{
  "mcpServers": {
    "bananamcp": {
      "command": "/path/to/bananamcp",
      "args": ["--vertex", "mcp"],
      "env": {
        "GCLOUD_PROJECT_ID": "your-project-id"
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
# OpenRouter (default)
bananamcp test A maine coon cat on a fancy throne holding a beer.

# Vertex AI
bananamcp --vertex test A futuristic banana spaceship.

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
├── main.go                    # Entry point (mcp + test subcommands, --vertex flag)
├── version.go                 # Semantic version
├── internal/
│   └── generator/
│       ├── generator.go       # ImageGenerator interface, factory, shared utilities
│       ├── openrouter.go      # OpenRouter backend (easy-llm-wrapper)
│       └── vertex.go          # Vertex AI backend (google.golang.org/genai)
├── Makefile
└── go.mod
```

The `ImageGenerator` interface in `internal/generator` abstracts backend selection. The factory function `generator.New(useVertex)` returns the appropriate implementation based on the `--vertex` flag. Both backends implement the same interface, so `main.go` is backend-agnostic.

## License

MIT
