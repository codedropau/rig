FROM golang:1.13 as build
WORKDIR /go/src/github.com/codedropau/rig
ADD . /go/src/github.com/codedropau/rig
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/rig-router github.com/codedropau/rig/cmd/rig-router

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/codedropau/rig/bin/rig-router /usr/local/bin/rig-router
