# InstagramRobot

<!-- [START badges] -->
<p>
  <!-- [GitHub Build Workflow] -->
  <a href="https://github.com/omegaatt36/instagramrobot/actions/workflows/build.yml"><img src="https://github.com/omegaatt36/instagramrobot/actions/workflows/build.yml/badge.svg" alt="Build status"></a>
  <!-- [GitHub Lint Workflow] -->
  <a href="https://github.com/omegaatt36/instagramrobot/actions/workflows/lint.yml"><img src="https://github.com/omegaatt36/instagramrobot/actions/workflows/lint.yml/badge.svg" alt="Build status"></a>
  <!-- [CodeFactor grade] -->
  <a href="https://codefactor.io/repository/github/omegaatt36/instagramrobot"><img src="https://www.codefactor.io/repository/github/omegaatt36/instagramrobot/badge" alt="CodeFactor"></a>
  <!-- [Go report score] -->
  <a href="https://goreportcard.com/report/github.com/omegaatt36/instagramrobot"><img src="https://goreportcard.com/badge/github.com/omegaatt36/instagramrobot?" alt="Go report" /></a>
  <!-- [GitHub license] -->
  <a href="https://github.com/omegaatt36/instagramrobot/blob/main/LICENSE"><img src="https://img.shields.io/github/license/omegaatt36/instagramrobot?color=blue" alt="License" /></a>
  <!-- [PRs welcome] -->
  <a href="https://github.com/omegaatt36/instagramrobot/pulls"><img src="https://img.shields.io/badge/PRs-welcome-blue.svg?color=d9ecde" alt="Contributing"></a>
</p>
<!-- [END badges] -->

<!-- [START description] -->

<a href="https://github.com/omegaatt36/instagramrobot" >
  <img align="right" src="https://raw.githubusercontent.com/omegaatt36/instagramrobot/main/images/ig-logo.svg" width="80" />
  <img align="right" src="https://raw.githubusercontent.com/omegaatt36/instagramrobot/main/images/telegram-logo.svg" width="80" />
  <img align="right" src="https://raw.githubusercontent.com/omegaatt36/instagramrobot/main/images/golang-logo.svg" height="80" />
</a>

> [InstagramRobot](https://github.com/omegaatt36/instagramrobot) is a bot based on [Telegram Bot API](https://core.telegram.org/bots/api) written in [Golang](https://golang.org/) that allows users to download public [Instagram](https://www.instagram.com/) photos, videos, and albums, without getting the user's credentials.

<!-- [END description] -->

## Table of contents

-   [Installing](#installing)
-   [Configuration](#configuration)
-   [Installing via Kubernetes by using helm](#installing-via-kubernetes-by-using-helm)
-   [Installing via Docker](#installing-via-docker)
    -   [Building the container](#building-the-container)
    -   [Running the container](#running-the-container)
-   [Installing as a service](#installing-as-a-service)
    -   [Build the application](#build-the-application)
    -   [Register the service](#register-the-service)
    -   [Enable the service at boot](#enable-the-service-at-boot)
    -   [Start the service](#start-the-service)


## Installing

<!--

You can download the latest version by checking the [GitHub releases](https://github.com/omegaatt36/instagramrobot/releases) page.

-->

Alternatively, you can download this project by cloning its Git repository:

```
git clone https://github.com/omegaatt36/instagramrobot.git
```

## Configuration

```bash
go run main.go --bot-token=***** --app-env=development
```

## Installing via Kubernetes by using helm

```shell
pushd charts
helm upgrade --install -v values.yaml insta-fetcher .
popd
```

## Installing via Docker

Docker is a tool designed to make it easier to create, deploy, and run applications by using containers.

Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.

If you're not familiar with Docker, [this guide](https://docs.docker.com/get-started/) is a great point to start.

### Building the container

```sh
docker compose build
```

### Running the container

```sh
docker compose up
```

## Installing as a service

Make sure that the project files exists in the `/usr/local/instagramrobot` directory.

### Build the application

> If you don't have Go installed, [click here](https://golang.org/doc/install) and follow its instructions.

```sh
go build -o bin/insta-fetcher
```

### Register the service

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

### Enable the service at boot

```sh
systemctl enable insta-fetcher
```

### Start the service

```sh
systemctl start insta-fetcher
```
