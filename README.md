# Telegram AI Bot

## Overview

This project is a feature-rich Telegram bot built with Go that integrates with cutting-edge AI APIs including OpenAI and
DeepSeek. The bot provides intelligent conversational capabilities, leveraging these powerful language models to deliver
high-quality responses to users.

## Prerequisites

- Docker and Docker Compose installed
- Go (if developing locally outside containers)
- Make utility (optional, but recommended for easier commands)

## Getting Started

### Development Environment

1. Clone the repository:
   ```bash
   git clone git@github.com:2m3j/telegram-ai-bot.git
   cd telegram-ai-bot
   ```
2. Copy the .env.dist to .env and update the environment variables:
   ```bash
   cp .env.dist .env
   ```
3. Edit .env with your credentials:
   ```dotenv
   TELEGRAM_TOKEN=your_bot_token
   AI_OPENAI_TOKEN=your_openai_token
   AI_DEEP_SEEK_TOKEN=your_deepseek_token
   ```
4. Start the development containers:
   ```bash
   make up
   ```
5. Run database migrations:
   ```bash
   make migrations-up
   ```
6. Connect to the debugger
   ```bash
   make connect
   ```
7. Once connected, you can use standard Delve commands.  

### Production Environment

1. To start production containers:
   ```bash
   make up-prod
   ```
2. Run production migrations:
   ```bash
   make migrations-up-prod
   ```

## Usage

![](/assets/demo.webp)

### Basic Commands

- `/start` - Welcome message and brief introduction
- `/ai` - Switch between AI providers
- `/clear` - Reset conversation context (bot will forget previous messages)

## Make Commands

### Development Commands

- `make up` - Start development containers
- `make connect` - Connect to the debugger
- `make stop` - Stop development containers
- `make down` - Stop and remove development containers
- `make docker-rm-volume` - Remove the MySQL volume
- `make migrations-up` - Apply all database migrations
- `make migrations-down` - Rollback all database migrations
- `make migrations-force-v version=X` - Force a specific migration version
- `make migrations-create name=description` - Create new migration files
- `make mock` - Generate mocks for interfaces
- `make test` - Run tests

### Production Commands

- `make up-prod` - Start production containers
- `make stop-prod` - Stop production containers
- `make down-prod` - Stop and remove production containers
- `make migrations-up-prod` - Apply production migrations
- `make migrations-down-prod` - Rollback production migrations
- `make migrations-force-v-prod version=X` - Force production migration version