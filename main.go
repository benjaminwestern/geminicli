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
	"strings"
)

// Documentation URL: https://ai.google.dev/api/rest/v1beta/SafetySetting#HarmBlockThreshold
const (
	baseURL         = "https://generativelanguage.googleapis.com"
	operatingSystem = runtime.GOOS
	githubURL       = "https://github.com/benjaminwestern/geminicli"
)

// API can count the number of tokens in the context and the conversation history
// curl https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:countTokens?key=$API_KEY \
//     -H 'Content-Type: application/json' \
//     -X POST \
//     -d '{
//       "contents": [{
//         "parts":[{
//           "text": "Write a story about a magic backpack."}]}]}' > response.json

func main() {
	conversationHistory := []Content{}
	contextFile := flag.String("context", "", "Path to the context file")
	outputPath := flag.String("output", "", "Path to the output file")
	debug := flag.Bool("debug", false, "Enable debug mode")
	hideWelcome := flag.Bool("hide-welcome", false, "Hide welcome message")
	tokenLimit := flag.Int("token-limit", 30720, "Max tokens for the conversation history")
	tokenWarning := flag.Int("token-warning", 25000, "Warning for the conversation history")
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

	outputTokenLimit := parseInt(importEnvironmentVariables("MAX_OUTPUT_TOKENS", "2048", *debug), *debug)
	temperature := parseFloat64(importEnvironmentVariables("TEMPERATURE", "0.9", *debug), *debug)
	topK := parseInt(importEnvironmentVariables("TOP_K", "1", *debug), *debug)
	topP := parseFloat64(importEnvironmentVariables("TOP_P", "1", *debug), *debug)

	context, err := loadContextFile(*contextFile)
	if err != nil {
		log.Fatal("Cannot load context file:", err)
	}

	err = validateInputContextTokenLimit(*tokenLimit, *tokenWarning, outputTokenLimit, context)
	if err != nil {
		fmt.Println("Context is too large...")
		fmt.Printf("Please enter the path to the context file (Max tokens %d or ~%d words):", *tokenLimit, *tokenLimit*5)
		var newContext string
		fmt.Scanln(&newContext)
		context, err = loadContextFile(newContext)
		if err != nil {
			log.Fatal("Cannot load context file:", err)
		}
		fmt.Println("Context changed.")
	}

	reader := bufio.NewReader(os.Stdin)

	config := GenerationConfig{
		Temperature:     temperature,
		TopK:            topK,
		TopP:            topP,
		MaxOutputTokens: outputTokenLimit,
		StopSequences:   []any{},
	}

	harassment := validateSafetyThreshold("HARASSMENT", *debug)
	hateSpeech := validateSafetyThreshold("HATE_SPEECH", *debug)
	sexuallyExplicit := validateSafetyThreshold("SEXUALLY_EXPLICIT", *debug)
	dangerousContent := validateSafetyThreshold("DANGEROUS_CONTENT", *debug)

	safetySettings := CreateDefaultSafetySettings(harassment, hateSpeech, sexuallyExplicit, dangerousContent)
	firstInput := true

	modelType := importEnvironmentVariables("MODEL_TYPE", "gemini-1.0-pro", *debug)
	apiVersion := importEnvironmentVariables("API_VERSION", "v1beta", *debug)

	url := fmt.Sprintf("%s/%s/models/%s:generateContent", baseURL, apiVersion, modelType)

	geminiURL := fmt.Sprintf("%s?key=%s", url, token)

	logFile, err := createTimestampedLogFile(*outputPath)
	if err != nil {
		log.Fatal("Cannot create markdown file:", err)
	}
	defer logFile.Close()

	if !*hideWelcome {
		fmt.Println("\n----------------- Welcome! -------------------")
		fmt.Println("Forget the browser. Chat with Gemini Pro right here!\n")
		fmt.Println("Add context to the conversation by adding a context file with the -context flag.\n")
		fmt.Printf("Your conversation will be logged in a markdown file. here %v\n", logFile.Name())
		fmt.Println("It might take a few seconds to get a response from the model. Please be patient.\n")
		fmt.Printf("Checkout the readme for more information on how to use this program at %v\n", githubURL)
	}

	for {
		fmt.Println("Enter your message or type 'menu' to see the options:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "menu" {
			displayMenu()
			var choice string
			fmt.Scanln(&choice)

			switch choice {

			case "1":
				clearScreen()
				fmt.Println("Resetting chat...")
				conversationHistory = []Content{}
				firstInput = true
				logFile, err = createTimestampedLogFile(*outputPath)
				if err != nil {
					log.Fatal("Cannot create markdown file:", err)
				}
				defer logFile.Close()
				fmt.Println("Chat reset.")

			case "2":
				clearScreen()
				fmt.Println("Change context...")
				fmt.Println("Enter the path to the context file:")
				var newContext string
				fmt.Scanln(&newContext)
				context, err = loadContextFile(newContext)
				if err != nil {
					log.Fatal("Cannot load context file:", err)
				}
				fmt.Println("Context changed.")
				conversationHistory = []Content{}
				firstInput = true

			case "3":
				clearScreen()
				fmt.Println("Delete conversation log...")
				files, err := os.ReadDir(".")
				if err != nil {
					log.Fatal("Cannot read directory:", err)
				}
				fmt.Println("Note: Only markdown files with the prefix 'conversation' will be displayed and you can't delete the current conversation log.")
				fmt.Println("Conversation logs:")
				for _, file := range files {
					if strings.HasSuffix(file.Name(), ".md") && strings.HasPrefix(file.Name(), "conversation") {
						if file.Name() == logFile.Name() {
							fmt.Println(file.Name())
						}
					}
				}
				fmt.Println("Enter the name of the file you want to delete:")
				var fileToDelete string
				fmt.Scanln(&fileToDelete)
				err = os.Remove(fileToDelete)
				if err != nil {
					log.Fatal("Cannot delete file:", err)
				}
				fmt.Println("Conversation log deleted.")

			case "4":
				clearScreen()
				fmt.Println("Change API Key...")
				fmt.Println("Enter the new API Key:")
				var newAPIKey string
				fmt.Scanln(&newAPIKey)
				token = newAPIKey
				geminiURL = fmt.Sprintf("%s?key=%s", url, token)
				fmt.Println("API Key changed.")

			case "5":
				clearScreen()
				fmt.Println("Exiting...")
				os.Exit(0)
			}

		} else {

			if firstInput && context != "" {
				input = context + "\n" + input
			}

			firstInput = false

			userInput := CreateContent("user", input)
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

			modelResponse := CreateContent("model", output.Candidates[0].Content.Parts[0].Text)

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

			err = validateConversationTokenLimit(*tokenLimit, *tokenWarning, conversationHistory)

			if err != nil {
				fmt.Println(err)
				convLength := len(conversationHistory)
				conversationHistory = conversationHistory[convLength-2:]
				fmt.Println("The last two messages have been removed from the conversation history.")
			}
		}
	}
}
