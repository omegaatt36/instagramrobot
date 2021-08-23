# InstagramRobot

<!-- [START badges] -->

<a href="https://github.com/feelthecode/instagramrobot/actions/workflows/build.yml"><img src="https://github.com/feelthecode/instagramrobot/actions/workflows/build.yml/badge.svg" alt="CI status" /></a>
<a href="https://goreportcard.com/report/github.com/feelthecode/instagramrobot"><img src="https://goreportcard.com/badge/github.com/feelthecode/instagramrobot" alt="Go report" /></a>
<a href="https://www.codefactor.io/repository/github/feelthecode/instagramrobot"><img src="https://www.codefactor.io/repository/github/feelthecode/instagramrobot/badge" alt="CodeFactor" /></a>

<!-- [END badges] -->

<!-- [START description] -->

<a href="https://github.com/feelthecode/instagramrobot">
  <img src="https://raw.githubusercontent.com/feelthecode/instagramrobot/main/images/ig-logo.svg" width="60" align="right" />
</a>
<a href="https://github.com/feelthecode/instagramrobot">
  <img src="https://raw.githubusercontent.com/feelthecode/instagramrobot/main/images/telegram-logo.svg" width="60" align="right" />
</a>

> [InstagramRobot](https://github.com/feelthecode/instagramrobot) is a bot based on [Telegram Bot API](https://core.telegram.org/bots/api) written in Golang that allows users to download Instagram photos, videos, and albums without providing their credentials.

<!-- [END description] -->

## Table of contents

-   [Installing as a service](#installing-as-a-service)
    -   [Creating configuration file](#creating-configuration-file)
    -   [Register service](#register-service)
    -   [Enable service at boot](#enable-service-at-boot)
    -   [Start service](#start-service)
-   [Docker](#docker)
    -   [Building container](#building-container)
    -   [Running container](#running-container)

## Installing as a service

First, make sure you're in the correct directory by executing the command below:

```
cd /usr/local/
```

You can download this project by cloning the Git repository:

```
git clone https://github.com/feelthecode/instagramrobot.git
cd instagramrobot
```

Alternatively, you can download the latest version by checking the [releases](https://github.com/feelthecode/instagramrobot/releases) page.

### Creating configuration file

```bash
# Create etc directory
mkdir etc

# Create .env file based on .env.example file
cp .env.example etc/.env

# Customize .env values
nano etc/.env
```

### Register service

```
nano /etc/systemd/system/igbot.service
```

```
[Unit]
Description=Telegram Instagram Bot Service

[Service]
WorkingDirectory=/usr/local/instagramrobot/bin
User=root
ExecStart=/usr/local/instagramrobot/bin/igbot --config-path /usr/local/instagramrobot/etc/
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target
```

### Enable service at boot

```
systemctl enable igbot
```

### Start service

```
systemctl start igbot
```

## Docker

Docker is a tool designed to make it easier to create, deploy, and run applications by using containers.

Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.

If you're not familiar with Docker, [this guide](https://docs.docker.com/get-started/) is a great point to start.

### Building container

```
docker-compose build
```

### Running container

```
docker-compose --env-file ./.env up
```
