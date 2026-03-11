package generator

import (
	"context"
	"fmt"
	"os"

	llm "github.com/nealhardesty/easy-llm-wrapper"
)

const defaultOpenRouterModel = "google/gemini-3.1-flash-image-preview"

// OpenRouterGenerator generates images via OpenRouter using easy-llm-wrapper.
type OpenRouterGenerator struct {
	client *llm.Client
}

func newOpenRouter() (*OpenRouterGenerator, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable is required")
	}

	model := resolveModel(defaultOpenRouterModel)

	client, err := llm.NewClientWithConfig(llm.Config{
		Provider: llm.ProviderOpenRouter,
		Model:    model,
		BaseURL:  "https://openrouter.ai/api/v1",
		APIKey:   apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	return &OpenRouterGenerator{client: client}, nil
}

// GenerateImage generates an image from the given prompt via OpenRouter and
// saves it to savePath. Returns the final path (which may have an extension appended).
func (g *OpenRouterGenerator) GenerateImage(ctx context.Context, prompt, savePath string) (string, error) {
	resp, err := g.client.Complete(ctx, llm.Request{
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

	return saveImage(savePath, resp.Images[0].Data, resp.Images[0].MIMEType)
}
