#!/usr/bin/make -f

export CGO_ENABLED=0

VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-list -1 HEAD)

# Builds the project.
build:
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/rig_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static" -X github.com/codedropau/rig/cmd/rig/version.GitVersion=${VERSION} -X github.com/codedropau/rig/cmd/rig/version.GitCommit=${COMMIT}' \
	    github.com/codedropau/rig/cmd/rig

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

IMAGE=codedropau/rig

# Releases the project Docker Hub.
release-docker:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

.PHONY: *