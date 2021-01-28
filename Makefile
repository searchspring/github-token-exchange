include .env
export
get:
	go get 
	go get github.com/golangci/golangci-lint/cmd/golangci-lint # Run lint
.PHONY: get

build:
	go build -v ./...
.PHONY: build


test-ci:
	go test -covermode=count -coverprofile=coverage.out ./...
.PHONY: test

test:
	go test -v -coverprofile=cover.out ./...
.PHONY: test-ci

run:
	PORT=3000 go run main.go github.go
.PHONY: run

docker-build:
	DOCKER_CONTENT_TRUST=1 && docker build -f Dockerfile -t github-token-exchange .
.PHONY: docker-build

docker-run:
	@docker run -e PORT='3000' \
		-e GITHUB_REDIRECT_URL="${GITHUB_REDIRECT_URL}" \
		-e GITHUB_CLIENT_ID="${GITHUB_CLIENT_ID}" \
		-e GITHUB_CLIENT_SECRET="${GITHUB_CLIENT_SECRET}" \
		-p 3000:3000 \
		github-token-exchange:latest
.PHONY: docker-run

fix-readme:
	npx remark-cli README.md -o
.PHONY: fix-readme
