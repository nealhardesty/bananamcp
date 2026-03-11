PRD: BananaMCP — Image Generator MCP Server (Golang)

1. Overview
A lightweight Model Context Protocol (MCP) server written in Go. It communicates via stdio and exposes a single tool to AI coding assistants (like Cursor or Claude). The tool allows the assistant to generate images and save them directly to the local filesystem. Two image generation backends are supported: OpenRouter (default) and Google Vertex AI (via the `--vertex` flag).

2. Tech Stack & Authentication
Language: Go (1.24+)

Protocol: Model Context Protocol (MCP) over stdio. Uses github.com/mark3labs/mcp-go.

LLM SDKs:
- OpenRouter (default): github.com/nealhardesty/easy-llm-wrapper — a unified Go LLM client. Image generation is performed via OpenRouter using NewClientWithConfig() with an explicit ProviderOpenRouter config.
- Vertex AI (--vertex): google.golang.org/genai — Google's official Gen AI Go SDK. Uses the Vertex AI backend with Application Default Credentials.

Authentication:
- OpenRouter: OPENROUTER_API_KEY is REQUIRED. The generator constructor fails fast with a clear error message if this variable is not set.
- Vertex AI: GCLOUD_PROJECT_ID is REQUIRED. Authentication uses Application Default Credentials (`gcloud auth application-default login` must have been run).

3. Environment Variables

OpenRouter mode (default):
- OPENROUTER_API_KEY (REQUIRED): OpenRouter API key.
- MODEL (optional): Override the default image generation model. Default: google/gemini-3.1-flash-image-preview

Vertex AI mode (--vertex):
- GCLOUD_PROJECT_ID (REQUIRED): Google Cloud project ID.
- GCLOUD_LOCATION (optional): Google Cloud region. Default: us-central1
- MODEL (optional): Override the default image generation model. Default: gemini-2.5-flash-image

Model priority: MODEL env var > built-in default.

4. Backend Abstraction

The ImageGenerator interface abstracts backend selection:

    type ImageGenerator interface {
        GenerateImage(ctx context.Context, prompt, savePath string) (string, error)
    }

A factory function New(useVertex bool) returns the appropriate implementation. Both OpenRouterGenerator and VertexGenerator implement this interface. The main package is backend-agnostic.

5. Tool Specification
Expose a single tool with the following schema:

Tool Name: generate_image

Description: "Generates an image based on a descriptive prompt and saves it to the specified local file path."

Parameters:

prompt (string, required): The detailed visual description and stylistic constraints for the image.

save_path (string, required): The base path where the image file should be saved. The actual file extension is derived from the MIME type returned by the model (e.g. .png, .jpg, .webp). If save_path already has an extension it is used as-is; otherwise the correct extension is appended.

6. Core Business Logic (Execution Flow)
When the generate_image tool is called, the server must:

Parse Inputs: Extract prompt and save_path from the tool request.

Backend Dispatch: The pre-configured ImageGenerator (selected at startup via --vertex flag) handles the generation.

OpenRouter path:
- Build an explicit OpenRouter config via easy-llm-wrapper
- Call client.Complete() with the prompt
- Extract image data from resp.Images[0]

Vertex AI path:
- Call client.Models.GenerateContent() with ResponseModalities set to IMAGE
- Extract image data from resp.Candidates[0].Content.Parts (InlineData)

File System Prep: Ensure the directory structure for save_path exists (os.MkdirAll).

Derive Extension: Use the MIME type to determine the correct file extension. Append to save_path if no extension is already present.

Save to Disk: Write image data to the final path.

Return:
- Success: "Successfully generated and saved image to: <final_path>"
- Error: Clear, human-readable error string for all failure modes.

7. Out of Scope
Hosting an HTTP/SSE server (strict stdio only).

Ollama or any non-OpenRouter/Vertex AI provider.

Image editing, in-painting, or multi-image composition.

State management or database integrations.

8. Tools
Always add a Makefile.

Add a simple command line test tool (test subcommand). It uses the same backend abstraction and default model, and saves the generated image with an extension derived from the response MIME type. Usage: bananamcp [--vertex] test <prompt>. Output defaults to output.<ext>.
