package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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

func CreateDefaultSafetySettings(harassment, hateSpeech, sexuallyExplicit, dangerousContent string) []SafetySettings {
	outputSafetySettings := []SafetySettings{}
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_HARASSMENT",
		Threshold: harassment,
	})
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_HATE_SPEECH",
		Threshold: hateSpeech,
	})
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
		Threshold: sexuallyExplicit,
	})
	outputSafetySettings = append(outputSafetySettings, SafetySettings{
		Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
		Threshold: dangerousContent,
	})
	return outputSafetySettings
}

func loadContextFile(filename string) (string, error) {
	if filename == "" {
		return "", nil // No context file specified
	}

	contextData, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading context file: %w", err)
	}
	return string(contextData), nil
}

func clearScreen() {
	var cmd *exec.Cmd
	if operatingSystem == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error clearing screen:", err)
	}
}

func validateSafetyThreshold(safetySetting string) string {
	currentSafetySetting := importEnvironmentVariables(safetySetting, "BLOCK_NONE")
	upperCast := strings.ToUpper(currentSafetySetting)

	// Allowed list pulled from the documentation for the API
	// https://ai.google.dev/api/rest/v1beta/SafetySetting#HarmBlockThreshold
	safetyAllowedList := []string{
		"HARM_BLOCK_THRESHOLD_UNSPECIFIED",
		"BLOCK_LOW_AND_ABOVE",
		"BLOCK_MEDIUM_AND_ABOVE",
		"BLOCK_ONLY_HIGH",
		"BLOCK_NONE",
	}

	for _, allowedSetting := range safetyAllowedList {
		if upperCast == allowedSetting {
			return upperCast
		}
	}

	fmt.Printf("Invalid safety setting: %s, using default: BLOCK_NONE\n", upperCast)
	return "BLOCK_NONE"
}

func importEnvironmentVariables(variable, defaultVar string) string {
	value := os.Getenv(variable)
	if value == "" {
		value = defaultVar
		fmt.Printf("Environment variable %s is not set, using defaults %s\n", variable, defaultVar)
	}
	return value
}

func parseFloat64(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Printf("Failed to parse float64, using default value: 0.9 | error: %v", err)
		return 0.9
	}
	return result
}

func parseInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Failed to parse int, using default value: 1 | error: %v", err)
		return 1
	}
	return result
}

func createTimestampedLogFile(path string) (*os.File, error) {
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	filename := "conversation-" + timestamp + ".md"
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get current working directory:", err)
		os.Exit(1)
	}

	// If the user has specified a path, use it
	if path != "" {
		filename = path + "/" + filename
	} else {
		filename = cwd + "/" + filename
	}

	logFile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}
