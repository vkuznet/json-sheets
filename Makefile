APP     = jsonsheets
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-s -w -extldflags '-static'"

.PHONY: build clean

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(APP) .

clean:
	rm -f $(APP)
