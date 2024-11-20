
# Goatmeal

A terminal-based chat application using the Groq API.

## Features

- Interactive TUI using Bubble Tea
- Markdown rendering support
- Conversation history
- Customizable themes
- System prompt management
- SQLite storage for chat history

## Installation

```bash
go install github.com/tedfulk/goatmeal@latest
```

## Configuration

On first run, Goatmeal will prompt you to:

1. Enter your Groq API key
2. Select a default model
3. Configure system prompt
4. Choose a theme

Configuration is stored in `~/.goatmeal/config.yaml`

## Usage

Simply run:

```bash
goatmeal
```
