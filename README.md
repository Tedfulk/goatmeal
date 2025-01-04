# Goatmeal - Terminal AI Chat & Web Search

Goatmeal is a powerful terminal-based application that provides access to various AI chat providers and web search capabilities, all within your terminal.

## Features

- **Multiple AI Providers Support**
  - OpenAI
  - Anthropic
  - Gemini
  - Deepseek
  - Groq
- **Web Search Integration**
  - Tavily search with domain filtering
  - Markdown-formatted search results
  - Answer summaries for relevant queries
- **User-Friendly Terminal UI**
  - Built with Bubble Tea and Bubbles
  - Beautiful styling with Lipgloss
  - Markdown rendering with Glamour
  - Multiple theme options
- **Conversation Management**
  - SQLite-based conversation storage
  - 30-day retention policy
  - Easy conversation browsing
  - Support for both chat and search conversations
- **Configuration**
  - YAML-based configuration
  - Secure API key storage
  - Customizable system prompts
  - Model selection per provider
- **Help System**
  - Built-in keyboard shortcut reference
  - Quick access with ctrl+h

## Installation

```bash
go install github.com/tedfulk/goatmeal@v1.1.16
```

## Configuration

On first run, Goatmeal will guide you through the setup process. You'll need to provide API keys for the services you want to use.

Configuration is stored in `~/.config/goatmeal/config.yaml`:

```yaml
api_keys:
  openai: "your-api-key"
  anthropic: "your-api-key"
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
  theme:
    name: "Default"
```

## Usage

### Keyboard Shortcuts

- `ctrl+t`: Start a new conversation
- `ctrl+l`: View conversation list
- `ctrl+s`: Open settings menu
- `ctrl+h`: View help
- `ctrl+c`: Quit application
- `esc`: Go back/close current view

### Chat Interface

- `?`: Toggle menu
- `/`: Enter search mode
- `enter`: Send message
- `#n`: Open message number 'n' in editor

### Search Mode

- `/query`: Search for information
- `/query +domain.com`: Search with specific domain
- `esc`: Exit search mode

### Conversation List

- `tab`: Switch focus between list and messages
- `d`: Delete selected conversation
- `esc`: Return to chat

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