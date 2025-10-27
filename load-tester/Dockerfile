# Start from golang:1.25-alpine base image
FROM golang:1.25-alpine AS build-env

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/*.go

FROM alpine:latest

# The latest alpine images don't have some tools like (`git` and `bash`).
#Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
     apk add --no-cache bash git openssh

WORKDIR /app

COPY --from=build-env /build/main main
COPY --from=build-env /build/config.json config.json

EXPOSE 4001

# Run the executable
ENTRYPOINT [ "./main", "start" ]
