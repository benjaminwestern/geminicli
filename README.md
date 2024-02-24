Title: Gemini AI Chatbot

Overview

This Go project builds a command-line chatbot powered by Google's Gemini 1.0 Pro large language model. It provides a conversational interface, safety settings, and the ability to log conversation transcripts to timestamped Markdown files.

Features

AI-Powered Chat: Interact with the state-of-the-art Gemini 1.0 Pro model for engaging and informative conversations.
Customizable: Adjust generation parameters like temperature and top-k/top-p for tailored responses.
Safety Considerations: Implement safety settings to help guide the model's output.
Conversation History: Track your chats with automatically generated Markdown logs, including timestamps.
New Chat Sessions: Easily start new chat sessions, with the option to clear previous conversation history.
Cross-Platform Support: Works on different operating systems thanks to dynamic screen clearing.
Prerequisites

A Google Cloud Platform project
A Gemini 1.0 Pro API key (set as the API_KEY environment variable)
Go programming language installed on your system
Getting Started

Clone this repository:

Bash
git clone https://github.com/benjaminwestern/gemini-ai-chatbot
Use code with caution.
Set your API key as an environment variable:

Bash
export API_KEY=your_gemini_api_key
Use code with caution.
Run the chatbot:

Bash
go run main.go
Use code with caution.
Usage

Enter your message in the terminal and press Enter.
The chatbot will process your message and generate a response using the Gemini AI model.
Type 'exit' to end the chat session.
If you want to start a new chat session, you have the option to remove the conversation history of the previous chat session.
Customization

Modify the following parameters within the main.go file to fine-tune the chatbot's behavior:

temperature: Controls the randomness of responses (higher values lead to more varied output).
topK: Limits the number of tokens considered for response generation.
topP: Controls the probability distribution for selecting the next token.
maxOutputTokens: Sets the maximum length of the generated response.
stopSequences: Define sequences that trigger the chatbot to stop generating text.
Safety Settings

Review and adjust the CreateDefaultSafetySettings() function to manage the AI's output in alignment with your desired safety guidelines. Refer to the Google AI API documentation for detailed safety settings customization: https://ai.google.dev/api/rest/v1beta/SafetySetting#HarmBlockThreshold

Contributing

We welcome contributions! Feel free to open issues for bug reports or feature requests. If you'd like to submit changes, please create a pull request.

Disclaimer

Please use this chatbot responsibly.  Remember that large language models can sometimes generate inaccurate, biased, or potentially harmful responses.

Let me know if you'd like any refinements or additional sections added to this README!
