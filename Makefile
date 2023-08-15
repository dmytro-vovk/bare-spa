.DEFAULT_GOAL := all

clean:
	@rm webserver 2> /dev/null || true

lint-front:
	@npm install && npm run lint

build-front:
	scripts/build-front.sh

ci-build-front:
	scripts/ci-build-front.sh

lint-back:
	@golangci-lint run

test-back:
	@go test -v ./internal/...

build-back:
	scripts/build-back.sh

assemble: build-front build-back

run: build-front
	@go run cmd/main.go

all: lint-front build-front lint-back test-back build-back clean
