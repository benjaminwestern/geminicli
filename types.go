package main

type GeminiRequest struct {
	Content          []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
	SafetySettings   []SafetySettings `json:"safetySettings"`
}

type Parts struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Parts `json:"parts"`
	Role  string  `json:"role"`
}
type GenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	TopK            int     `json:"topK"`
	TopP            int     `json:"topP"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
	StopSequences   []any   `json:"stopSequences"`
}

type SafetySettings struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiResponse struct {
	Candidates     []Candidates   `json:"candidates"`
	PromptFeedback PromptFeedback `json:"promptFeedback"`
}

type SafetyRatings struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}
type Candidates struct {
	Content       Content         `json:"content"`
	FinishReason  string          `json:"finishReason"`
	Index         int             `json:"index"`
	SafetyRatings []SafetyRatings `json:"safetyRatings"`
}
type PromptFeedback struct {
	SafetyRatings []SafetyRatings `json:"safetyRatings"`
}
