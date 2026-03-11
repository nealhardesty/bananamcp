package generator

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

const defaultVertexModel = "gemini-2.5-flash-image"

// VertexGenerator generates images via Google Vertex AI using the genai SDK.
type VertexGenerator struct {
	client *genai.Client
	model  string
}

func newVertex() (*VertexGenerator, error) {
	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("GCLOUD_PROJECT_ID environment variable is required for --vertex mode")
	}

	location := os.Getenv("GCLOUD_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	model := resolveModel(defaultVertexModel)

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %w", err)
	}

	return &VertexGenerator{client: client, model: model}, nil
}

// GenerateImage generates an image from the given prompt via Vertex AI and
// saves it to savePath. Returns the final path (which may have an extension appended).
func (g *VertexGenerator) GenerateImage(ctx context.Context, prompt, savePath string) (string, error) {
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
			Role: "user",
		},
	}

	resp, err := g.client.Models.GenerateContent(ctx, g.model, contents, &genai.GenerateContentConfig{
		ResponseModalities: []string{string(genai.ModalityImage)},
	})
	if err != nil {
		return "", fmt.Errorf("image generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("model returned no candidates")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			return saveImage(savePath, part.InlineData.Data, part.InlineData.MIMEType)
		}
	}

	return "", fmt.Errorf("model returned no images")
}
