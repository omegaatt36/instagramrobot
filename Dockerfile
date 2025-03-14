FROM golang:1.24-alpine as build

WORKDIR /go/src/app

COPY ["go.mod", "go.sum", "./"]

RUN ["go", "mod", "download"]

COPY . .

ENV CGO_ENABLED=0

ENV GOPROXY=https://proxy.golang.org

ENV APP_NAME=insta-fetcher

RUN ["go", "build", "-o", "build/${APP_NAME}", "cmd/bot/main.go"]

FROM build as dev

CMD ["go", "run", "."]

FROM gcr.io/distroless/static-debian12 as prod

WORKDIR /home/app/

COPY --from=build /go/src/app/build/${APP_NAME} ./

CMD ["./${APP_NAME}"]
