include .env
export

test:
	go test -v -coverprofile=cover.out ./...
.PHONY: test

run:
	PORT=3000 go run main.go github.go
.PHONY: run

docker-build:
	DOCKER_CONTENT_TRUST=1 && docker build -f Dockerfile -t github-token-exchange .
.PHONY: docker-build

docker-run:
	@echo running on port 3000
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
