FROM golang:1.24-bookworm as build

WORKDIR /go/src/app

COPY ["go.mod", "go.sum", "./"]

RUN ["go", "mod", "download"]

COPY . .

ENV GOPROXY=https://proxy.golang.org

ENV CGO_ENABLED=0

RUN ["go", "build", "-o", "insta-fetcher", "cmd/bot/main.go"]

FROM build as dev

CMD ["go", "run", "."]

FROM gcr.io/distroless/static-debian12 as prod

WORKDIR /home/app/

COPY --from=build /go/src/app/insta-fetcher ./

CMD ["./insta-fetcher"]
