# Goatmeal - Terminal AI Chat & Web Search

Goatmeal is a powerful terminal-based application that provides access to various AI chat providers and web search capabilities, all within your terminal.

https://github.com/user-attachments/assets/01bf6cda-39a6-4b41-97c1-fb7321bcd291

## Features

- **Multiple AI Providers Support**
  - OpenAI
  - Anthropic
  - Gemini
  - Deepseek
  - Groq
  - Ollama (local models)
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
go install github.com/tedfulk/goatmeal@v1.2.17
```

## Configuration

On first run, Goatmeal will guide you through the setup process. You'll need to provide API keys for the services you want to use.

### Ollama Setup

To use Ollama with Goatmeal:

1. Install Ollama from [ollama.ai](https://ollama.ai)
2. Start the Ollama server locally
3. No API key is required - Goatmeal will automatically connect to Ollama at `http://localhost:11434`
4. Select "ollama" as your provider in Goatmeal's settings to see available models

Configuration is stored in `~/.config/goatmeal/config.yaml`:

```yaml
api_keys:
  openai: your-api-key
  anthropic: your-api-key
  gemini: your-api-key
  deepseek: your-api-key
  groq: your-api-key
  tavily: your-api-key
  ollama: ollama # goatmeal will put a default api key in
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
- `/web query`: Search for information
- `/web query +domain.com`: Search with specific domain
- `/webe query`: Enhanced web search with AI optimization
- `/webe query +domain.com`: Enhanced domain-specific search
- `/epq`: Enhanced Programming query
- `enter`: Send message
- `/o[n]`: Open message number 'n' in editor (e.g., /o1)
- `/c[n]`: Copy message number 'n' to clipboard (e.g., /c1)
- `/b[n]`: Copy code block number 'n' to clipboard (e.g., /b1)
- `/s[n]`: Speak message number 'n' using system TTS (e.g., /s1)
- `ctrl+q`: Stop current speech playback

#### Enhanced Search

The enhanced search mode (🔍+) uses AI to optimize your search queries for better results. When using `/webe`:

Examples:

- Basic: `/web what's the latest news in the quantum computing?`
- Enhanced: `/webe what's the latest news in the quantum computing?` gets transformed into something like `Recent breakthroughs in quantum computing 2024-2025 including, advancements in quantum processors algorithms and applications from reputable sources like research, journals and tech news.`
  - Gets transformed into a more specific query including location and time context
- Domain-specific: `/webe python tutorials +python.org`
  - Enhanced query limited to python.org domain

#### Enhanced Programming query

The enhanced programming query mode (💻+) uses AI to optimize programming-related queries by adding specificity about languages, frameworks, and technical requirements.

Examples:

- Basic: `/epq how to build a web app with react and go`
- Enhanced: `What is a step-by-step guide to building a scalable and efficient web application using React as the frontend framework and Go as the backend server? Consider a RESTful API architecture, focusing on best practices for setting up a new React project with Create React App, designing a robust API with Go's net/http package or a framework like Gin, and implementing authentication and authorization using JSON Web Tokens (JWT) or OAuth.`

### Conversation List

- `tab`: Switch focus between list and messages
- `ctrl+d`: Delete selected conversation
- `ctrl+e`: Export conversation as JSON (saves to ~/Downloads)
- `esc`: Return to chat

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

### Text Selection in Terminal

When using Goatmeal, there are two ways to select and copy text:

1. **Using Option Key (Recommended)**:

   - Hold the Option (⌥) key while selecting text
   - This temporarily disables mouse reporting without affecting app functionality
   - Release Option to restore normal mouse interaction

2. **Toggle Mouse Reporting**:
   - Press ⌘⇧M to temporarily disable mouse reporting
   - Select text normally
   - Press ⌘⇧M again to re-enable mouse reporting

Note: Mouse reporting is required for proper UI interaction (scrolling, clicking, etc).
Completely disabling it will impair application functionality.
