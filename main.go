package main

import (
	"bytes"
	"encoding/json"
	"log"
	"fmt"
	"io"
	"time"
	"bufio"
	"net/http"
	"runtime"
	"os/exec"
	"os"
	"flag"
)

// Documentation URL: https://ai.google.dev/api/rest/v1beta/SafetySetting#HarmBlockThreshold
const (
	url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.0-pro:generateContent"
	operatingSystem = runtime.GOOS 
	githubURL = "https://github.com/benjaminwestern/geminicli"
)

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

func main() {
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
		Temperature:     0.9,
		TopK:            1,
		TopP:            1,
		MaxOutputTokens: 2048,
		StopSequences:   []any{},
	}
	safetySettings := CreateDefaultSafetySettings()
	firstInput := true

	geminiURL := fmt.Sprintf("%s?key=%s", url, token)

	conversationHistory := []Content{}

	logFile, err := createTimestampedLogFile(*outputPath)
	if err != nil {
		log.Fatal("Cannot create markdown file:", err)
	}
	defer logFile.Close()

	fmt.Println("Forget the browser. Chat with Gemini Pro right here!\n")
	fmt.Println("Add context to the conversation by adding a context file with the -context flag.\n")
	fmt.Printf("Your conversation will be logged in a markdown file. here %v/%v\n",logFile.Name())
	fmt.Println("It might take a few seconds to get a response from the model. Please be patient.\n")
	fmt.Printf("Checkout the readme for more information on how to use this program at %v\n", githubURL)

	for { 
		fmt.Println("Enter your message (or type 'exit' to quit):")
		input, _ := reader.ReadString('\n')	
		
		if input == "exit\n" { 
			break
		}

		if firstInput && context != ""{
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
		if err != nil {
			log.Fatal("failed to marshall gemini request:", err)
		}

		req, err := http.NewRequest("POST", geminiURL, bytes.NewBuffer(marshalledRequest))
		if err != nil {
			log.Fatal("failed to create http request:", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("failed to send http request:", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		output := GeminiResponse{}
		err = json.Unmarshal(body, &output)
		if err != nil {
			log.Fatal("failed to unmarshal gemini response:", err)
		}

		modelResponse := CreateModelContent(output.Candidates[0].Content.Parts[0].Text)

		fmt.Println("Model Response:", modelResponse.Parts[0].Text)
		finishReason := output.Candidates[0].FinishReason
		if finishReason != "STOP" {
			fmt.Println("Unknown finish Reason:", output.Candidates[0].FinishReason)
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

