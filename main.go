// Command bananamcp is an image generation tool and MCP server.
//
// Usage:
//
//	bananamcp [--vertex] mcp              Start the MCP stdio server (for use with AI assistants)
//	bananamcp [--vertex] test <prompt>    Generate an image and save to output.<ext>
//	bananamcp --version                   Print version and exit
//
// Examples:
//
//	bananamcp mcp
//	bananamcp --vertex mcp
//	bananamcp test A maine coon cat on a fancy throne holding a beer.
//	bananamcp --vertex test A futuristic banana spaceship.
//
// When --vertex is provided, Vertex AI is used instead of OpenRouter.
// OPENROUTER_API_KEY must be set for the default backend.
// GCLOUD_PROJECT_ID must be set for --vertex mode.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nealhardesty/bananamcp/internal/generator"
)

func main() {
	args := os.Args[1:]

	var useVertex bool
	if len(args) > 0 && args[0] == "--vertex" {
		useVertex = true
		args = args[1:]
	}

	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "mcp":
		gen, err := generator.New(useVertex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := runMCP(gen); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "test":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "error: prompt required\n\nUsage: bananamcp [--vertex] test <prompt>\n")
			os.Exit(1)
		}
		gen, err := generator.New(useVertex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := runTest(gen, strings.Join(args[1:], " ")); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "--version", "-v", "version":
		fmt.Printf("bananamcp %s\n", Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %q\n\n", args[0])
		usage()
		os.Exit(1)
	}
}

// runMCP starts the MCP stdio server using the provided image generator.
func runMCP(gen generator.ImageGenerator) error {
	s := server.NewMCPServer(
		"bananamcp",
		Version,
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("generate_image",
			mcp.WithDescription("Generates an image based on a descriptive prompt and saves it to the specified local file path."),
			mcp.WithString("prompt",
				mcp.Description("The detailed visual description and stylistic constraints for the image."),
				mcp.Required(),
			),
			mcp.WithString("save_path",
				mcp.Description("The base path where the image file should be saved. The actual file extension is derived from the MIME type returned by the model (e.g. .png, .jpg, .webp). If save_path already has an extension it is used as-is; otherwise the correct extension is appended."),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGenerateImage(ctx, req, gen)
		},
	)

	return server.ServeStdio(s)
}

// handleGenerateImage handles the generate_image MCP tool call.
func handleGenerateImage(ctx context.Context, req mcp.CallToolRequest, gen generator.ImageGenerator) (*mcp.CallToolResult, error) {
	prompt, err := req.RequireString("prompt")
	if err != nil {
		return mcp.NewToolResultError("prompt is required"), nil
	}

	savePath, err := req.RequireString("save_path")
	if err != nil {
		return mcp.NewToolResultError("save_path is required"), nil
	}

	finalPath, err := gen.GenerateImage(ctx, prompt, savePath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully generated and saved image to: %s", finalPath)), nil
}

// runTest generates an image from prompt and saves it to output.<ext>.
func runTest(gen generator.ImageGenerator, prompt string) error {
	finalPath, err := gen.GenerateImage(context.Background(), prompt, "output")
	if err != nil {
		return err
	}
	fmt.Printf("Successfully generated and saved image to: %s\n", finalPath)
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: bananamcp [--vertex] <subcommand> [args]

Subcommands:
  mcp              Start the MCP stdio server
  test <prompt>    Generate an image and save to output.<ext>
  version          Print version and exit

Flags:
  --vertex         Use Google Vertex AI instead of OpenRouter

Environment variables (OpenRouter mode - default):
  OPENROUTER_API_KEY  (required) OpenRouter API key
  MODEL               (optional) Override default model (google/gemini-3.1-flash-image-preview)

Environment variables (Vertex AI mode - --vertex):
  GCLOUD_PROJECT_ID   (required) Google Cloud project ID
  GCLOUD_LOCATION     (optional) Google Cloud region (default: us-central1)
  MODEL               (optional) Override default model (gemini-2.5-flash-image)

Examples:
  bananamcp mcp
  bananamcp --vertex mcp
  bananamcp test A maine coon cat on a fancy throne holding a beer.
  bananamcp --vertex test A futuristic banana spaceship.
`)
}
