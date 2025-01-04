# Goatmeal - Terminal AI Chat & Web Search

Goatmeal is a powerful terminal-based application that provides access to various AI chat providers and web search capabilities, all within your terminal.

## Features

- **Multiple AI Providers Support**
  - OpenAI (GPT-4, GPT-3.5)
  - Claude
  - Gemini
  - Deepseek
  - Groq
- **Web Search Integration**
  - Tavily search integration
- **User-Friendly Terminal UI**
  - Built with Bubble Tea and Bubbles
  - Beautiful styling with Lipgloss
  - Markdown rendering with Glamour
- **Conversation Management**
  - SQLite-based conversation storage
  - 30-day retention policy
  - Easy conversation browsing
- **Configuration**
  - YAML-based configuration
  - Secure API key storage
  - Customizable system prompts
  - Model selection per provider

## Installation

```bash
go install github.com/tedfulk/goatmeal@v1.1.15
```

## Configuration

On first run, Goatmeal will guide you through the setup process. You'll need to provide API keys for the services you want to use.

Configuration is stored in `~/.config/goatmeal/config.yaml`:

```yaml
api_keys:
  openai: "your-api-key"
  claude: "your-api-key"
  gemini: "your-api-key"
  deepseek: "your-api-key"
  groq: "your-api-key"
  tavily: "your-api-key"

system_prompts:
  default: "You are a helpful assistant."
  code_helper: "You are an expert code assistant."
  creative: "You are a creative assistant."

settings:
  output_glamour: true
  default_models:
    openai: "gpt-4-turbo-preview"
    claude: "claude-3-opus"
    gemini: "gemini-pro"
    deepseek: "deepseek-coder"
    groq: "mixtral-8x7b-32768"
```

## Usage

### Basic Commands

- Start the application:
  ```bash
  goatmeal
  ```

### Keyboard Shortcuts

- `Ctrl+C` or `q`: Quit
- `Ctrl+S`: Open settings
- `Enter`: Send message
- `Up/Down`: Navigate conversation history
- `Tab`: Switch between input and conversation view

## Development

### Prerequisites

- Go 1.23 or higher
- SQLite

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/tedfulk/goatmeal.git
   cd goatmeal
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build:
   ```bash
   go build
   ```

### Project Structure

```
goatmeal/
├── config/           # Configuration management
├── chat/
│   └── providers/    # AI provider implementations
├── search/          # Web search integration
├── database/        # SQLite database management
├── ui/             # Terminal UI components
└── main.go         # Application entry point
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details. 