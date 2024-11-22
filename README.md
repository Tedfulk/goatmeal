
# Goatmeal

A terminal-based chat application using the Groq API.

## Features

- Interactive TUI using Bubble Tea
- Markdown rendering support
- Conversation history
- Customizable themes
- System prompt management
- SQLite storage for chat history
- Image input

## Installation

```bash
go install github.com/tedfulk/goatmeal@v1.1.0
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
