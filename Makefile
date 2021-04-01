export APP=virtual-box
export LDFLAGS="-w -s"

run-server:
	go run -ldflags $(LDFLAGS) ./cmd/virtual-box server

build:
	go build -ldflags $(LDFLAGS) ./cmd/virtual-box

install:
	go install -ldflags $(LDFLAGS) ./cmd/virtual-box
