package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

// Documentation URL: https://ai.google.dev/api/rest/v1beta/SafetySetting#HarmBlockThreshold
const (
	baseURL         = "https://generativelanguage.googleapis.com"
	operatingSystem = runtime.GOOS
	githubURL       = "https://github.com/benjaminwestern/geminicli"
)

func main() {
	conversationHistory := []Content{}
	contextFile := flag.String("context", "", "Path to the context file")
	outputPath := flag.String("output", "", "Path to the output file")
	help := flag.Bool("help", false, "Show help message")

	flag.Usage = func() {
		fmt.Println("Usage: geminicli [flags]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	token := os.Getenv("API_KEY")
	if token == "" {
		log.Fatal("API_KEY is not set")
		os.Exit(1)
	}

	context, err := loadContextFile(*contextFile)
	if err != nil {
		log.Fatal("Cannot load context file:", err)
	}

	reader := bufio.NewReader(os.Stdin)

	config := GenerationConfig{
		Temperature:     parseFloat64(importEnvironmentVariables("TEMPERATURE", "0.9")),
		TopK:            parseInt(importEnvironmentVariables("TOP_K", "1")),
		TopP:            parseInt(importEnvironmentVariables("TOP_P", "1")),
		MaxOutputTokens: parseInt(importEnvironmentVariables("MAX_OUTPUT_TOKENS", "2048")),
		StopSequences:   []any{},
	}

	harassment := validateSafetyThreshold("HARASSMENT")
	hateSpeech := validateSafetyThreshold("HATE_SPEECH")
	sexuallyExplicit := validateSafetyThreshold("SEXUALLY_EXPLICIT")
	dangerousContent := validateSafetyThreshold("DANGEROUS_CONTENT")

	safetySettings := CreateDefaultSafetySettings(harassment, hateSpeech, sexuallyExplicit, dangerousContent)
	firstInput := true

	modelType := importEnvironmentVariables("MODEL_TYPE", "gemini-1.0-pro")
	apiVersion := importEnvironmentVariables("API_VERSION", "v1beta")

	url := fmt.Sprintf("%s/%s/models/%s:generateContent", baseURL, apiVersion, modelType)

	geminiURL := fmt.Sprintf("%s?key=%s", url, token)

	logFile, err := createTimestampedLogFile(*outputPath)
	if err != nil {
		log.Fatal("Cannot create markdown file:", err)
	}
	defer logFile.Close()

	fmt.Println("\n----------------- Welcome! -------------------")
	fmt.Println("Forget the browser. Chat with Gemini Pro right here!\n")
	fmt.Println("Add context to the conversation by adding a context file with the -context flag.\n")
	fmt.Printf("Your conversation will be logged in a markdown file. here %v\n", logFile.Name())
	fmt.Println("It might take a few seconds to get a response from the model. Please be patient.\n")
	fmt.Printf("Checkout the readme for more information on how to use this program at %v\n", githubURL)

	for {
		fmt.Println("Enter your message (or type 'exit' to quit):")
		input, _ := reader.ReadString('\n')

		if input == "exit\n" {
			break
		}

		if firstInput && context != "" {
			input = context + "\n" + input
		}

		firstInput = false
		userInput := CreateUserContent(input)
		conversationHistory = append(conversationHistory, userInput)

		geminiRequest := GeminiRequest{
			GenerationConfig: config,
			SafetySettings:   safetySettings,
			Content:          conversationHistory,
		}

		marshalledRequest, err := json.Marshal(geminiRequest)

		fmt.Println("Sending request to Gemini...")
		if err != nil {
			log.Fatal("failed to marshall gemini request:", err)
			os.Exit(1)
		}

		req, err := http.NewRequest("POST", geminiURL, bytes.NewBuffer(marshalledRequest))
		if err != nil {
			log.Fatal("failed to create http request:", err)
			os.Exit(1)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("failed to send http request:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Fatal("failed to get 200 status code from gemini:", err)
			os.Exit(1)
		}

		fmt.Println("Received response from Gemini...")
		body, _ := io.ReadAll(resp.Body)
		output := GeminiResponse{}
		err = json.Unmarshal(body, &output)
		if err != nil {
			log.Fatal("failed to unmarshal gemini response:", err)
			os.Exit(1)
		}

		modelResponse := CreateModelContent(output.Candidates[0].Content.Parts[0].Text)

		fmt.Println("Model Response:", modelResponse.Parts[0].Text)
		finishReason := output.Candidates[0].FinishReason
		if finishReason != "STOP" {
			// TODO Handle these:
			// https://ai.google.dev/api/rest/v1beta/Candidate#finishreason
			// FINISH_REASON_UNSPECIFIED	Default value. This value is unused.
			// STOP	Natural stop point of the model or provided stop sequence.
			// MAX_TOKENS	The maximum number of tokens as specified in the request was reached.
			// SAFETY	The candidate content was flagged for safety reasons.
			// RECITATION	The candidate content was flagged for recitation reasons.
			// OTHER	Unknown reason.
			fmt.Println("Unhandled finish Reason:", output.Candidates[0].FinishReason)
		}

		conversationHistory = append(conversationHistory, modelResponse)

		fmt.Fprintf(logFile, "**User:** %s\n", userInput.Parts[0].Text)
		fmt.Fprintf(logFile, "**Model:** %s\n\n", modelResponse.Parts[0].Text)

		fmt.Println("Are you finished with this chat? (yes/no)")
		var startNewChat string
		fmt.Scanln(&startNewChat)

		if startNewChat == "yes" {
			fmt.Println("Remove the conversation history of the last chat? (yes/no)")
			var removeHistory string
			fmt.Scanln(&removeHistory)

			if removeHistory == "yes" {
				os.Remove(logFile.Name())
			}

			fmt.Println("Exit the program? (yes/no)")
			var exitProgram string
			fmt.Scanln(&exitProgram)

			if exitProgram == "yes" {
				clearScreen()
				os.Exit(0)
			}

			firstInput = true
			conversationHistory = []Content{}

			logFile, err = createTimestampedLogFile(*outputPath)
			if err != nil {
				log.Fatal("Cannot create markdown file:", err)
			}
			defer logFile.Close()

			clearScreen()
			fmt.Println("----------------- New chat started -------------------")

		} else if startNewChat == "exit" {
			break
		}
	}
}
