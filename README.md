# InstagramRobot


[![wakatime](https://wakatime.com/badge/user/61550660-cb83-43cc-9bc5-cf742c36b4cd/project/5ee6b14b-f44a-4330-a9ce-ff24385d9c28.svg)](https://wakatime.com/badge/user/61550660-cb83-43cc-9bc5-cf742c36b4cd/project/5ee6b14b-f44a-4330-a9ce-ff24385d9c28)
[![Docker Build and Push](https://github.com/omegaatt36/instagramrobot/actions/workflows/release.yml/badge.svg)](https://github.com/omegaatt36/instagramrobot/actions/workflows/release.yml)
[![Build status](https://github.com/omegaatt36/instagramrobot/actions/workflows/lint.yml/badge.svg)](https://github.com/omegaatt36/instagramrobot/actions/workflows/lint.yml)
[![CodeFactor](https://www.codefactor.io/repository/github/omegaatt36/instagramrobot/badge)](https://codefactor.io/repository/github/omegaatt36/instagramrobot)
[![Go report](https://goreportcard.com/badge/github.com/omegaatt36/instagramrobot)](https://goreportcard.com/report/github.com/omegaatt36/instagramrobot)
[![License](https://img.shields.io/github/license/omegaatt36/instagramrobot?color=blue)](https://github.com/omegaatt36/instagramrobot/blob/main/LICENSE)
[![Contributing](https://img.shields.io/badge/PRs-welcome-blue.svg?color=d9ecde)](https://github.com/omegaatt36/instagramrobot/pulls)

InstagramRobot is a Golang project that allows you to download public Instagram and Threads content without requiring login credentials. It supports:

- Instagram posts (single image/video)
- Instagram Reels
- Instagram multi-image/video posts (Carousel/Album)
- Threads posts (text, images, videos, multimedia)

This project offers two main ways to use it:

1.  Telegram Bot: Send an Instagram or Threads link directly to the bot, and it will return the downloaded media files and text content.
2.  Web UI: A simple web interface based on HTMX. Paste an Instagram link, and you can preview the media and text content directly in your browser. (Note: The current Web UI only supports Instagram).

END description

## Table of contents

- [Usage](#usage)
  - [Telegram Bot](#telegram-bot)
  - [Web UI](#web-ui)
- [Deployment Options](#deployment-options)
  - [Running Directly](#running-directly)
  - [Using Docker](#using-docker)
  - [Using Kubernetes (Helm)](#using-kubernetes-helm)
  - [As a Systemd Service](#as-a-systemd-service)
- [Development](#development)

## Usage

### Telegram Bot

How it works:
Send a public Instagram or Threads post/Reels link directly to your deployed Telegram Bot.
The bot will automatically parse the link, download the media (images and videos), and send them back to the chat along with the post's caption.
It supports both single media items and albums (multiple media).

How to deploy your own Bot:
Refer to the "Deployment Options" section below and choose the method that suits you. You will need a Telegram Bot Token.

### Web UI

Note: The current Web UI only supports Instagram links, not Threads.

How to access (running locally):

```shell
go run cmd/web/main.go
```

Then open http://localhost:8080 in your browser.

How to use:
Paste a public Instagram post/Reels link into the input field and click "Submit".
The page will dynamically load and display the post's media (images/videos) and caption.
You can click on images or videos to enlarge them for preview.

How to deploy the Web UI:
Refer to the "Deployment Options" section. The web service listens on port 8080 by default.

## Deployment Options

You can deploy the Telegram Bot or the Web UI using any of the following methods, depending on your environment.

### Running Directly

This is the most basic way to run the application, suitable for local testing or simple deployments. You need to have Go installed first.

```bash
# Clone the repository (if you haven't already)
git clone https://github.com/omegaatt36/instagramrobot.git
cd instagramrobot

# Run the Telegram Bot
# Replace YOUR_BOT_TOKEN with your actual Telegram Bot token
go run cmd/bot/main.go --bot-token=YOUR_BOT_TOKEN --app-env=development

# Run the Web UI
go run cmd/web/main.go --app-env=development

# You can use --app-env=production for production settings (e.g., JSON logs)
# You can also configure via environment variables, for example:
# export BOT_TOKEN=YOUR_BOT_TOKEN
# export APP_ENV=production
# go run cmd/bot/main.go
```

### Using Docker

This is the recommended deployment method, using containers to package the application and its dependencies.

Building the container image:

```sh
# Build for both bot and web (multi-stage build likely defined in Dockerfile)
docker compose build
# Or, if you build images separately or push to registry (like in Taskfile.yaml):
# docker buildx build --platform linux/amd64,linux/arm64 -t your-repo/insta-fetcher:latest . --push
```

Running the container (using Docker Compose):
Check the `deploy/docker-compose.yml` file to ensure environment variables (like `BOT_TOKEN`) are set.

```sh
# Make sure BOT_TOKEN is set in deploy/docker-compose.yml or your environment
docker compose -f deploy/docker-compose.yml up -d # Use -d to run in detached mode
```

This Compose file might start both the Bot and Web services, or you might need to choose. Please check its contents.

### Using Kubernetes (Helm)

If you use Kubernetes, you can deploy using the provided Helm chart.

```shell
# Navigate to the chart directory
pushd deploy/charts/bot # Assuming the chart is primarily for the bot, adjust if needed

# Install or upgrade the release named 'insta-fetcher'
# You'll likely need to configure values like botToken via values.yaml or --set
helm upgrade --install -f values.yaml insta-fetcher .

# Go back to the previous directory
popd
```

You will need to edit `deploy/charts/bot/values.yaml` to set your Bot Token and other configurations.

### As a Systemd Service

Suitable for deployment on traditional Linux servers.

Build the application binary:

```sh
# Build the bot binary
go build -o bin/insta-fetcher cmd/bot/main.go
# Build the web binary if needed
# go build -o bin/web-server cmd/web/main.go
```

Assuming you place the project files in `/usr/local/instagramrobot`.

Register the service:
Create the `/etc/systemd/system/insta-fetcher.service` file (example for the Bot):

```ini
[Unit]
Description=Telegram Instagram Bot Service
After=network.target

[Service]
# User=your_service_user # It's better to run as a non-root user
# Group=your_service_group
WorkingDirectory=/usr/local/instagramrobot
# Set environment variables if preferred over command-line args
# Environment="BOT_TOKEN=YOUR_BOT_TOKEN"
# Environment="APP_ENV=production"
# Environment="LOG_LEVEL=info"
ExecStart=/usr/local/instagramrobot/bin/insta-fetcher --bot-token=YOUR_BOT_TOKEN --app-env=production --log-level=info
Restart=on-failure
RestartPreventExitStatus=23 # Optional: prevent restart on specific exit codes
StandardOutput=journal # Log to systemd journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Remember to replace `YOUR_BOT_TOKEN`. Consider creating a dedicated non-root user to run the service. You can use either `--bot-token` or `Environment="BOT_TOKEN=..."`.

Enable the service at boot:

```sh
sudo systemctl enable insta-fetcher.service
```

Start the service:

```sh
sudo systemctl start insta-fetcher.service
```

Check service status:

```sh
sudo systemctl status insta-fetcher.service
# View logs
sudo journalctl -u insta-fetcher.service -f
```

## Development

This project uses `go-task` (Taskfile.yaml) to manage the development workflow.

Install development dependencies:

```sh
task dependency
```

Format code:

```sh
task fmt
```

Run linters and checks:

```sh
task check
```

Run tests:

```sh
task test
```

Live reload for Web UI development:
Requires `air` to be installed first: `go install github.com/cosmtrek/air@latest`

```sh
task live-web
```

This monitors Go file changes, automatically recompiling and restarting the web server.
