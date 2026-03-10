// Package generator provides image generation logic using the easy-llm-wrapper library via OpenRouter.
package generator

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	llm "github.com/nealhardesty/easy-llm-wrapper"
)

const defaultImageModel = "google/gemini-3.1-flash-image-preview"

// GenerateImage generates an image from the given prompt and saves it to savePath.
// The actual file extension is derived from the MIME type returned by the model.
// If savePath already has an extension it is used as-is; otherwise the correct extension is appended.
// Returns the final path where the image was saved.
func GenerateImage(ctx context.Context, prompt, savePath string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable is not set")
	}

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
	if err != nil {
		return "", fmt.Errorf("failed to create LLM client: %w", err)
	}

	resp, err := client.Complete(ctx, llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleUser, Parts: []llm.Part{llm.TextPart(prompt)}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("image generation failed: %w", err)
	}

	if len(resp.Images) == 0 {
		return "", fmt.Errorf("model returned no images")
	}

	finalPath := derivePath(savePath, resp.Images[0].MIMEType)

	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(finalPath, resp.Images[0].Data, 0644); err != nil {
		return "", fmt.Errorf("failed to write image: %w", err)
	}

	return finalPath, nil
}

// derivePath returns the final save path. If savePath already has an extension
// it is returned unchanged; otherwise the extension derived from mimeType is appended.
func derivePath(savePath, mimeType string) string {
	if filepath.Ext(savePath) != "" {
		return savePath
	}
	return savePath + extFromMIME(mimeType)
}

// extFromMIME returns the file extension (including leading dot) for a MIME type.
func extFromMIME(mimeType string) string {
	mediaType, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		mediaType = mimeType
	}

	switch strings.ToLower(mediaType) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		exts, err := mime.ExtensionsByType(mediaType)
		if err == nil && len(exts) > 0 {
			return exts[0]
		}
		return ".bin"
	}
}
