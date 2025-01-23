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
go install github.com/tedfulk/goatmeal@v1.2.6
```

## Configuration

On first run, Goatmeal will guide you through the setup process. You'll need to provide API keys for the services you want to use.

Configuration is stored in `~/.config/goatmeal/config.yaml`:

```yaml
api_keys:
  openai: your-api-key
  anthropic: your-api-key
  gemini: your-api-key
  deepseek: your-api-key
  groq: your-api-key
  tavily: your-api-key
current_model: llama-3.3-70b-versatile
current_provider: groq
current_system_prompt: You are a helpful AI assistant.
settings:
    outputglamour: true
    conversationretention: 30
    theme:
        name: Default
    username: teddy
system_prompts:
    - content: You are a helpful AI assistant.
      title: General
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
- `#o[n]`: Open message number 'n' in editor (e.g., #o1)
- `#c[n]`: Copy message number 'n' to clipboard (e.g., #c1)
- `#b[n]`: Copy code block number 'n' to clipboard (e.g., #b1)
- `#s[n]`: Speak message number 'n' using system TTS (e.g., #s1)
- `ctrl+k`: Stop current speech playback

### Search Mode

- `/query`: Search for information
- `/query +domain.com`: Search with specific domain
- `/o`: Toggle enhanced search mode
- `/o query`: One-time enhanced search
- `esc`: Exit search mode

#### Enhanced Search

The enhanced search mode (🔍+) uses AI to optimize your search queries for better results. When enabled:

- Queries are automatically enhanced with relevant context
- Location and time information are considered
- Ambiguous terms are clarified
- Results are more targeted and comprehensive
- Queries are optimized while staying under character limits

Examples:

- Basic: `/what's the weather like?`
- Enhanced: `/o what's the weather like?`
  - Gets transformed into a more specific query including location and time context
- Domain-specific: `/o python tutorials +python.org`
  - Enhanced query limited to python.org domain

The status bar shows 🔍+ when enhanced mode is active. You can toggle it with `/o` or use it for a single search with `/o query`.

### Conversation List

- `tab`: Switch focus between list and messages
- `d`: Delete selected conversation
- `e`: Export conversation as JSON (saves to ~/Downloads)
- `esc`: Return to chat

### Project Structure

```md
goatmeal/
 ├── config
 ├── database
 ├── main.go
 ├── scripts
 ├── services
 │   ├── providers
 │   ├── web
 ├── ui
 │   ├── setup
 │   ├── theme
 │   ├── // the rest of the views and logic
 └── utils
```

## Dependencies

### Linux

For text-to-speech functionality on Linux, you'll need to install `espeak`:

```bash
# Ubuntu/Debian
sudo apt-get install espeak

# Fedora
sudo dnf install espeak

# Arch Linux
sudo pacman -S espeak
```

### macOS

Text-to-speech is supported out of the box using the built-in `say` command.
