package main

func CreateUserContent(input string) Content {
	return Content{
		Role:  "user",
		Parts: []Parts{{Text: input}},
	}
}

func CreateModelContent(input string) Content {
	return Content{
		Role:  "model",
		Parts: []Parts{{Text: input}},
	}
}

func CreateDefaultSafetySettings() []SafetySettings {
	outputSafetySettings := []SafetySettings{}
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_HARASSMENT",
		Threshold: "BLOCK_NONE",
	})
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_HATE_SPEECH",
		Threshold: "BLOCK_NONE",
	})
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
		Threshold: "BLOCK_NONE",
	})
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
		Threshold: "BLOCK_NONE",
	})
	return outputSafetySettings
}
