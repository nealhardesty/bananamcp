// Package generator provides image generation backends for BananaMCP.
// It supports OpenRouter (via easy-llm-wrapper) and Google Vertex AI
// (via the google.golang.org/genai SDK).
package generator

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// ImageGenerator generates images from text prompts and saves them to disk.
type ImageGenerator interface {
	GenerateImage(ctx context.Context, prompt, savePath string) (string, error)
}

// New creates an ImageGenerator for the selected backend.
// When useVertex is true, the Vertex AI backend is used; otherwise OpenRouter.
func New(useVertex bool) (ImageGenerator, error) {
	if useVertex {
		return newVertex()
	}
	return newOpenRouter()
}

// saveImage derives the final file path, ensures the directory exists,
// and writes imageData to disk. It returns the final path.
func saveImage(savePath string, imageData []byte, mimeType string) (string, error) {
	finalPath := derivePath(savePath, mimeType)

	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(finalPath, imageData, 0644); err != nil {
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

// resolveModel returns the MODEL env var if set, otherwise the provided default.
func resolveModel(defaultModel string) string {
	if m := os.Getenv("MODEL"); m != "" {
		return m
	}
	return defaultModel
}
