FROM golang:1.19-alpine as build

# Set the working directory
WORKDIR /go/src/app

# Cache dependencies
COPY ["go.mod", "go.sum", "./"]

# Download dependencies
RUN ["go", "mod", "download"]

# Copy project files
COPY . .

# The cgo tool is enabled by default for native builds on systems where it is expected to work.
# It is disabled by default when cross-compiling
ENV CGO_ENABLED=0

# Controls the source of Go module downloads
# Can help assure builds are deterministic and secure.
ENV GOPROXY=https://proxy.golang.org

# Executable filename (binary file)
ENV APP_NAME=igbot

# Build binary file
RUN ["go", "build", "-o", "build/${APP_NAME}"]

#
# Development build
#
FROM build as dev

# Run the application via Go
CMD ["go", "run", "."]

#
# Production build
#
FROM alpine:3.14.1 as prod

# By default, Docker runs container as root which inside of the container can pose as a security issue.
RUN addgroup -S app && adduser -S -G app app
USER app

# Set the working directory
WORKDIR /home/app/

COPY --from=build /go/src/app/build/${APP_NAME} ./

# Execute the binary file
CMD ["./${APP_NAME}"]
