#!/usr/bin/env bash

# usage:
# ./run_tests.sh                         # local, go 1.11
# ./run_tests.sh docker                  # docker, go 1.11
# ./run_tests.sh podman                  # podman, go 1.11

set -ex

# The script does automatic checking on a Go package and its sub-packages,
# including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. go vet        (http://golang.org/cmd/vet)
# 3. gosimple      (https://github.com/dominikh/go-simple)
# 4. unconvert     (https://github.com/mdempsky/unconvert)
# 5. ineffassign   (https://github.com/gordonklaus/ineffassign)
# 6. race detector (http://blog.golang.org/race-detector)

# golangci-lint (github.com/golangci/golangci-lint) is used to run each each
# static checker.

# Default GOVERSION
[[ ! "$GOVERSION" ]] && GOVERSION=1.11
REPO=dcrextdata

testrepo () {
  export GO111MODULE=on

  go version
  

  # Test application install
  go build

  # Get linter
  go get github.com/golangci/golangci-lint/cmd/golangci-lint

  env GORACE='halt_on_error=1' go test -v -race ./...

  # check linters
  golangci-lint run --deadline=10m --disable-all --enable govet --enable staticcheck \
    --enable gosimple --enable unconvert --enable ineffassign --enable structcheck \
    --enable goimports --enable misspell --enable unparam


  echo "------------------------------------------"
  echo "Tests completed successfully!"

}

DOCKER=
[[ "$1" == "docker" || "$1" == "podman" ]] && DOCKER=$1
if [ ! "$DOCKER" ]; then
    testrepo
    exit
fi

# Don't really know what other image to use yet
DOCKER_IMAGE_TAG=dcrdata-golang-builder-$GOVERSION
$DOCKER pull decred/$DOCKER_IMAGE_TAG

$DOCKER run --rm -it -v $(pwd):/src decred/$DOCKER_IMAGE_TAG /bin/bash -c "\
  rsync -ra --include-from=<(git --git-dir=/src/.git ls-files) \
  --filter=':- .gitignore' \
  /src/ /go/src/github.com/raedahgroup/$REPO/ && \
  cd github.com/raedahgroup/$REPO/ && \
  env GOVERSION=$GOVERSION GO111MODULE=on bash run_tests.sh"
