# InstagramRobot

<!-- [START badges] -->
<p>
  <!-- [GitHub Build Workflow] -->
  <a href="https://github.com/feelthecode/instagramrobot/actions/workflows/build.yml"><img src="https://github.com/feelthecode/instagramrobot/actions/workflows/build.yml/badge.svg" alt="Build status"></a>
  <!-- [GitHub Lint Workflow] -->
  <a href="https://github.com/feelthecode/instagramrobot/actions/workflows/lint.yml"><img src="https://github.com/feelthecode/instagramrobot/actions/workflows/lint.yml/badge.svg" alt="Build status"></a>
  <!-- [CodeFactor grade] -->
  <a href="https://codefactor.io/repository/github/feelthecode/instagramrobot"><img src="https://www.codefactor.io/repository/github/feelthecode/instagramrobot/badge" alt="CodeFactor"></a>
  <!-- [Go report score] -->
  <a href="https://goreportcard.com/report/github.com/feelthecode/instagramrobot"><img src="https://goreportcard.com/badge/github.com/feelthecode/instagramrobot" alt="Go report" /></a>
  <!-- [GitHub license] -->
  <a href="https://github.com/feelthecode/instagramrobot/blob/main/LICENSE"><img src="https://img.shields.io/github/license/feelthecode/instagramrobot?color=blue" alt="License" /></a>
  <!-- [PRs welcome] -->
  <a href="https://github.com/feelthecode/instagramrobot/pulls"><img src="https://img.shields.io/badge/PRs-welcome-blue.svg?color=d9ecde" alt="Contributing"></a>
</p>
<!-- [END badges] -->

<!-- [START description] -->

<a href="https://github.com/feelthecode/instagramrobot" >
  <img align="right" src="https://raw.githubusercontent.com/feelthecode/instagramrobot/main/images/ig-logo.svg" width="80" />
  <img align="right" src="https://raw.githubusercontent.com/feelthecode/instagramrobot/main/images/telegram-logo.svg" width="80" />
  <img align="right" src="https://raw.githubusercontent.com/feelthecode/instagramrobot/main/images/golang-logo.svg" height="80" />
</a>

> [InstagramRobot](https://github.com/feelthecode/instagramrobot) is a bot based on [Telegram Bot API](https://core.telegram.org/bots/api) written in [Golang](https://golang.org/) that allows users to download public [Instagram](https://www.instagram.com/) photos, videos, and albums, without getting the user's credentials.

<!-- [END description] -->

## Table of contents

-   [Installing](#installing)
-   [Configuration](#configuration)
-   [Installing via Docker](#installing-via-docker)
    -   [Building the container](#building-the-container)
    -   [Running the container](#running-the-container)
-   [Installing as a service](#installing-as-a-service)
    -   [Build the application](#build-the-application)
    -   [Register the service](#register-the-service)
    -   [Enable the service at boot](#enable-the-service-at-boot)
    -   [Start the service](#start-the-service)

## Installing

You can download the latest version by checking the [GitHub releases](https://github.com/feelthecode/instagramrobot/releases) page.

Alternatively, you can download this project by cloning its Git repository:

```
git clone https://github.com/feelthecode/instagramrobot.git
```

## Configuration

Create the `.env` file based on the `.env.example` file, and update its values.

```bash
cp .env.example .env
```

## Installing via Docker

Docker is a tool designed to make it easier to create, deploy, and run applications by using containers.

Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.

If you're not familiar with Docker, [this guide](https://docs.docker.com/get-started/) is a great point to start.

### Building the container

```
docker-compose build
```

### Running the container

```
docker-compose up
```

## Installing as a service

Make sure that the project files exists in the `/usr/local/instagramrobot` directory.

### Build the application

> If you don't have Go installed, [click here](https://golang.org/doc/install) and follow its instructions.

```
go build -o bin/igbot
```

### Register the service

Start by creating the `/etc/systemd/system/igbot.service` file.

```
[Unit]
Description=Telegram Instagram Bot Service

[Service]
WorkingDirectory=/usr/local/instagramrobot/bin
User=root
ExecStart=/usr/local/instagramrobot/bin/igbot --config-path [CONFIG_PATH]
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target
```

Don't forget to replace the `[CONFIG_PATH]` with the correct path of configuration file.

### Enable the service at boot

```
systemctl enable igbot
```

### Start the service

```
systemctl start igbot
```
