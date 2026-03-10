PRD: BananaMCP — Image Generator MCP Server (Golang)

1. Overview
A lightweight Model Context Protocol (MCP) server written in Go. It communicates via stdio and exposes a single tool to AI coding assistants (like Cursor or Claude). The tool allows the assistant to generate images using the `easy-llm-wrapper` library via OpenRouter and save them directly to the local filesystem.

2. Tech Stack & Authentication
Language: Go (1.21+)

Protocol: Model Context Protocol (MCP) over stdio. Use a standard community Go MCP library (e.g., github.com/mark3labs/mcp-go).

LLM SDK: github.com/nealhardesty/easy-llm-wrapper — a unified Go LLM client. Image generation is performed exclusively via OpenRouter using NewClientWithConfig() with an explicit ProviderOpenRouter config. The generic NewClient() env-based auto-detection is NOT used because Ollama does not support image generation and the model must always be an image-generation model.

Authentication: OPENROUTER_API_KEY is REQUIRED. The server must fail fast with a clear error message if this variable is not set. There is no fallback provider.

3. Environment Variables
OPENROUTER_API_KEY (REQUIRED): OpenRouter API key. Server refuses to start without it.

MODEL (optional): Override the default image generation model. Default: google/gemini-3.1-flash-image-preview

Model priority (matching elwi convention): MODEL env var > built-in default.

The server must NOT read OLLAMA_HOST or any Vertex AI / GCP environment variables.

4. Tool Specification
Expose a single tool with the following schema:

Tool Name: generate_image

Description: "Generates an image based on a descriptive prompt and saves it to the specified local file path."

Parameters:

prompt (string, required): The detailed visual description and stylistic constraints for the image.

save_path (string, required): The base path where the image file should be saved. The actual file extension is derived from the MIME type returned by the model (e.g. .png, .jpg, .webp). If save_path already has an extension it is used as-is; otherwise the correct extension is appended.

5. Core Business Logic (Execution Flow)
When the generate_image tool is called, the server must:

Parse Inputs: Extract prompt and save_path from the tool request.

File System Prep: Ensure the directory structure for save_path exists (os.MkdirAll).

LLM Client: Fail immediately if OPENROUTER_API_KEY is not set. Build an explicit OpenRouter config (matching elwi):

    const defaultImageModel = "google/gemini-3.1-flash-image-preview"

    apiKey := os.Getenv("OPENROUTER_API_KEY")
    // error if empty

    model := defaultImageModel
    if m := os.Getenv("MODEL"); m != "" {
        model = m
    }

    client, err := llm.NewClientWithConfig(llm.Config{
        Provider: llm.ProviderOpenRouter,
        Model:    model,
        BaseURL:  "https://openrouter.ai/api/v1",
        APIKey:   apiKey,
    })

API Call: Call client.Complete(ctx, llm.Request{...}) with the prompt as the sole user message.

Process Response: Check resp.Images. If empty, return an error ("model returned no images").

Derive Extension: Use the MIME type of resp.Images[0] to determine the correct file extension (image/png → .png, image/jpeg → .jpg, image/webp → .webp, etc.). Append to save_path if no extension is already present.

Save to Disk: Write resp.Images[0].Data to the final path.

Return:

Success: "Successfully generated and saved image to: <final_path>"

Error: Clear, human-readable error string for all failure modes (missing API key, network error, no image in response, filesystem error).

6. Out of Scope
Hosting an HTTP/SSE server (strict stdio only).

Ollama or any non-OpenRouter provider.

Image editing, in-painting, or multi-image composition.

State management or database integrations.

7. Tools
Always add a Makefile.

Add a simple command line test tool (cmd/banana) modelled directly on elwi from easy-llm-wrapper. It uses the same explicit OpenRouter config, same default model (google/gemini-3.1-flash-image-preview), same model priority (MODEL env > built-in default), and saves the generated image with an extension derived from the response MIME type. Usage: ./banana <prompt> ex. ./banana A maine coon cat on a fancy throne holding a beer. Output defaults to output.<ext>. Requires OPENROUTER_API_KEY.
