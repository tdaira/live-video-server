FROM golang:1.13.1 AS build-env
ADD . $GOPATH/src/github.com/tdaira/live-video-server/
WORKDIR $GOPATH/src/github.com/tdaira/live-video-server/
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=$(git rev-parse --verify HEAD)" \
    -o live-video-server ./cmd/server/server.go

FROM alpine:3.10.2
COPY --from=build-env /go/src/github.com/tdaira/live-video-server/live-video-server /usr/local/bin/live-video-server
COPY ./config.toml .
COPY ./video ./video/
RUN apk add --no-cache tzdata ca-certificates
EXPOSE 80

CMD ["live-video-server"]
