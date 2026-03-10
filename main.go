// Command bananamcp is an image generation tool and MCP server.
//
// Usage:
//
//	bananamcp mcp              Start the MCP stdio server (for use with AI assistants)
//	bananamcp test <prompt>    Generate an image and save to output.<ext>
//	bananamcp --version        Print version and exit
//
// Examples:
//
//	bananamcp mcp
//	bananamcp test A maine coon cat on a fancy throne holding a beer.
//
// OPENROUTER_API_KEY must be set. Optionally set MODEL to override the default
// image generation model (google/gemini-3.1-flash-image-preview).
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
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "mcp":
		if err := runMCP(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "test":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "error: prompt required\n\nUsage: bananamcp test <prompt>\n")
			os.Exit(1)
		}
		if err := runTest(strings.Join(os.Args[2:], " ")); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "--version", "-v", "version":
		fmt.Printf("bananamcp %s\n", Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %q\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

// runMCP starts the MCP stdio server.
func runMCP() error {
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		return fmt.Errorf("OPENROUTER_API_KEY environment variable is required")
	}

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
		handleGenerateImage,
	)

	return server.ServeStdio(s)
}

// handleGenerateImage handles the generate_image MCP tool call.
func handleGenerateImage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	prompt, err := req.RequireString("prompt")
	if err != nil {
		return mcp.NewToolResultError("prompt is required"), nil
	}

	savePath, err := req.RequireString("save_path")
	if err != nil {
		return mcp.NewToolResultError("save_path is required"), nil
	}

	finalPath, err := generator.GenerateImage(ctx, prompt, savePath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully generated and saved image to: %s", finalPath)), nil
}

// runTest generates an image from prompt and saves it to output.<ext>.
func runTest(prompt string) error {
	finalPath, err := generator.GenerateImage(context.Background(), prompt, "output")
	if err != nil {
		return err
	}
	fmt.Printf("Successfully generated and saved image to: %s\n", finalPath)
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: bananamcp <subcommand> [args]

Subcommands:
  mcp              Start the MCP stdio server
  test <prompt>    Generate an image and save to output.<ext>
  version          Print version and exit

Environment variables:
  OPENROUTER_API_KEY  (required) OpenRouter API key
  MODEL               (optional) Override default model (google/gemini-3.1-flash-image-preview)

Examples:
  bananamcp mcp
  bananamcp test A maine coon cat on a fancy throne holding a beer.
`)
}
