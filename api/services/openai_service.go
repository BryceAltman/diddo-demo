package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"diddo-api/configs"
	"diddo-api/models"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAIService struct {
	client *openai.Client
	config *configs.Config
}

func NewOpenAIService(config *configs.Config) *OpenAIService {
	client := openai.NewClient(config.OpenAIKey)
	return &OpenAIService{
		client: client,
		config: config,
	}
}

func (s *OpenAIService) IdentifyClothingItems(imageData []byte) (models.ClothingIdentificationResponse, error) {
	if s.config.OpenAIKey == "" {
		return models.ClothingIdentificationResponse{}, fmt.Errorf("OpenAI API key not configured")
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	prompt := `Identify and analyze all clothing items in this image. For each of them return a string to be used to find them item with a google shopping search, be extremely descriptive so that an exact match can be found. Analyze the image carefully and identify all visible clothing items. Return only valid JSON.`

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o20240806,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: prompt,
						},
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL: fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
							},
						},
					},
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name: "clothing_item_descriptions",
					Schema: &jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"items": {
								Type: jsonschema.Array,
								Items: &jsonschema.Definition{
									Type: jsonschema.Object,
									Properties: map[string]jsonschema.Definition{
										"description": {
											Type: jsonschema.String,
										},
									},
									Required: []string{"description"},
								},
							},
						},
						Required: []string{"items"},
					},
				},
			},
		},
	)

	if err != nil {
		return models.ClothingIdentificationResponse{}, fmt.Errorf("failed to identify clothing items: %w", err)
	}

	if len(resp.Choices) == 0 {
		return models.ClothingIdentificationResponse{}, fmt.Errorf("no response from OpenAI")
	}

	var openAIResponse struct {
		Items []models.ClothingItemDescription `json:"items"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &openAIResponse); err != nil {
		return models.ClothingIdentificationResponse{}, fmt.Errorf("failed to parse OpenAI response: %w. Response: %s", err, resp.Choices[0].Message.Content)
	}

	return models.ClothingIdentificationResponse{
		Items:  openAIResponse.Items,
		Status: "success",
	}, nil
}
