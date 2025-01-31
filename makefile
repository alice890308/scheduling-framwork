
BIN_DIR=_output/bin

# If tag not explicitly set in users default to the git sha.
TAG ?= ${shell (git describe --tags --abbrev=14 | sed "s/-g\([0-9a-f]\{14\}\)$/+\1/") 2>/dev/null || git rev-parse --verify --short HEAD}

all: local

init:
	mkdir -p ${BIN_DIR}

local: init
	go build -o=${BIN_DIR}/alice-scheduler ./cmd/scheduler
	#go build -o=_output/bin/alice-scheduler ./cmd/scheduler

build-linux: init
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=${BIN_DIR}/alice-scheduler ./cmd/scheduler
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o=_output/bin/alice-scheduler ./cmd/scheduler

image: build-linux
	docker build --no-cache . -t alice-scheduler:${TAG}

clean:
	rm -rf _output/
	rm -f *.log
