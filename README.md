# InstagramRobot

<!-- [START badges] -->

[![wakatime](https://wakatime.com/badge/user/61550660-cb83-43cc-9bc5-cf742c36b4cd/project/5ee6b14b-f44a-4330-a9ce-ff24385d9c28.svg)](https://wakatime.com/badge/user/61550660-cb83-43cc-9bc5-cf742c36b4cd/project/5ee6b14b-f44a-4330-a9ce-ff24385d9c28)
[![Build status](https://github.com/omegaatt36/instagramrobot/actions/workflows/build.yml/badge.svg)](https://github.com/omegaatt36/instagramrobot/actions/workflows/build.yml)
[![Build status](https://github.com/omegaatt36/instagramrobot/actions/workflows/lint.yml/badge.svg)](https://github.com/omegaatt36/instagramrobot/actions/workflows/lint.yml)
[![CodeFactor](https://www.codefactor.io/repository/github/omegaatt36/instagramrobot/badge)](https://codefactor.io/repository/github/omegaatt36/instagramrobot)
[![Go report](https://goreportcard.com/badge/github.com/omegaatt36/instagramrobot)](https://goreportcard.com/report/github.com/omegaatt36/instagramrobot)
[![License](https://img.shields.io/github/license/omegaatt36/instagramrobot?color=blue)](https://github.com/omegaatt36/instagramrobot/blob/main/LICENSE)
[![Contributing](https://img.shields.io/badge/PRs-welcome-blue.svg?color=d9ecde)](https://github.com/omegaatt36/instagramrobot/pulls)
<!-- [END badges] -->

<!-- [START description] -->

> [InstagramRobot](https://github.com/omegaatt36/instagramrobot) is a bot based on [Telegram Bot API](https://core.telegram.org/bots/api) written in [Golang](https://golang.org/) that allows users to download public [Instagram](https://www.instagram.com/) photos, videos, and albums, without getting the user's credentials.

<!-- [END description] -->

## Table of contents

- [Installing](#installing-telegram-bot)
- [Configuration](#configuration)
- [Installing via Kubernetes by using helm](#installing-via-kubernetes-by-using-helm)
- [Installing via Docker](#installing-via-docker)
  - [Building the container](#building-the-container)
  - [Running the container](#running-the-container)
- [Installing as a service](#installing-as-a-service)
  - [Build the application](#build-the-application)
  - [Register the service](#register-the-service)
  - [Enable the service at boot](#enable-the-service-at-boot)
  - [Start the service](#start-the-service)

## Installing Telegram Bot

Alternatively, you can download this project by cloning its Git repository:

```bash
git clone https://github.com/omegaatt36/instagramrobot.git
```

### Configuration

```bash
go run main.go --bot-token=***** --app-env=development
```

### Installing via Kubernetes by using helm

```shell
pushd deploy/charts/bot
helm upgrade --install -v values.yaml insta-fetcher .
popd
```

### Installing via Docker

Docker is a tool designed to make it easier to create, deploy, and run applications by using containers.

Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.

If you're not familiar with Docker, [this guide](https://docs.docker.com/get-started/) is a great point to start.

#### Building the container

```sh
docker compose build
```

#### Running the container

```sh
docker compose up -f deploy/docker-compose.yml
```

### Installing as a service

Make sure that the project files exist in the `/usr/local/instagramrobot` directory.

#### Build the application

> If you don't have Go installed, [click here](https://golang.org/doc/install) and follow its instructions.

```sh
go build cmd/bot/main.go -o bin/insta-fetcher
```

#### Register the service

Start by creating the `/etc/systemd/system/insta-fetcher.service` file.

```sh
[Unit]
Description=Telegram Instagram Bot Service

[Service]
WorkingDirectory=/usr/local/instagramrobot/bin
User=root
ExecStart=/usr/local/instagramrobot/bin/insta-fetcher --bot-token=FILL_ME
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target
```

Don't forget to replace the `--bot-token=FILL_ME` with the correct bot token.

#### Enable the service at boot

```sh
systemctl enable insta-fetcher
```

#### Start the service

```sh
systemctl start insta-fetcher
```

## Start Web Version

Powered by [HTMX](https://htmx.org/)

```shell
go run cmd/web/main.go
```

Then you can open `http://localhost:8080` on your browser.
