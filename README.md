# GeminiCLI  

GeminiCLI is a command-line interface tool written in Go, designed to interact with the generative language API provided by Google. It enables users to generate content by sending requests to the API and receiving responses directly in the terminal. This tool is particularly useful for developers and researchers working on natural language processing and generative AI models.  
## Features  
- Send requests to the generative language API.
- Save conversation history in markdown files with timestamps.
- Clear screen functionality for better user experience. 
- Load context from a file to influence the generation. 
- Command-line flags for flexible configuration.  

## Pre-Requisites
1. Before installing GeminiCLI, ensure you have Go installed on your system. You can download and install Go from [the official Go website](https://golang.org/dl/).  
2. Setup an API Key for Gemini by following this documentation - [Gemini API Key](https://aistudio.google.com/app/apikey)

## Installation  

To install GeminiCLI, clone the repository to your local machine:  
```bash 
git clone https://github.com/benjaminwestern/geminicli.git cd geminicli
```

Change directory to the geminicli repo you just cloned, then, build the binary using:

```bash 
go build .
```

This will create the `geminicli` executable in your current directory.
## Usage
To start using GeminiCLI, you can run the executable directly from the command line. Here are some common commands:

- To generate content with a context file:

```bash
./geminicli -context path/to/your/context.txt
```

- To specify an output file for the conversation log:

```bash
./geminicli -output path/to/your/output/
```

- For help and a list of available flags:

```bash
./geminicli -help
```

Ensure you have set the `API_KEY` environment variable with your API key before running the tool.
## Contributing
I welcome contributions from the community! If you'd like to contribute to GeminiCLI, please fork the repository and submit a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.
## License
GeminiCLI is open-sourced under the MIT License. See the LICENSE file for more details.
## Getting Help
If you encounter any issues or have questions about using GeminiCLI, please open an issue on the GitHub issue tracker.

For more information about the generative language API and its configuration, visit [Google's API documentation](https://ai.google.dev/api/rest).
## Important Information
Please ensure you are aware that anything you input into these models (if your using a personal google account) will be utilised to further train the Gemini models, so PLEASE think twice before sending private content to these models!
[Google's Gemini FAQs](https://gemini.google.com/faq)
