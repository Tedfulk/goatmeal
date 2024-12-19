
# Goatmeal

A terminal-based chat application using the Groq API.

https://github.com/user-attachments/assets/a21b9fa0-7949-4521-8c34-b047c6d1c30d

## Features

- Interactive TUI using Bubble Tea
- Markdown rendering support (via [glamour](https://github.com/charmbracelet/glamour))
- Conversation history
- Themes (via [lipgloss](https://github.com/charmbracelet/lipgloss))
- System prompt management
- SQLite storage for chat history
- Image input (remote images only, couldn't figure out how to get local image size small enough to send to api)

## Installation

```bash
go install github.com/tedfulk/goatmeal@v1.1.12
```

## Configuration

On first run, Goatmeal will prompt you to:

1. Enter your Groq [API key](https://console.groq.com/keys)
2. Select a default model
3. Configure system prompt
4. Choose a theme
5. Enter your username

Configuration is stored in `~/.goatmeal/config.yaml`

## Usage

Simply run:

```bash
goatmeal
```

## Shortcuts

| Shortcut | Action |
| --- | --- |
| shift+tab | Toggle menu |
| tab | Toggle focus (in chat/list view) |
| esc | Back/Quit |
| q | Quit |
| enter | Send message |
| shift+enter | New line in message |
| ctrl+l | List conversations |
| ctrl+t | New conversation |
| ; | Go to theme selector |
| # | Go to image input |
| ↑/k | Scroll up |
| ↓/j | Scroll down |
| ↑/k | Previous item |
| ↓/j | Next item |
| enter | Select item |

## Roadmap

- [ ] Add support for ChatGPT
- [ ] Add support for conversation export
- [ ] Add support for file uploads
