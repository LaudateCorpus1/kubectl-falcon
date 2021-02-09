GO=GO111MODULE=on go
GOBUILD=$(GO) build

build:
	$(GOBUILD) ./...


GORELEASER_FLAGS    = -tags="$(GORELEASER_BUILD_TAGS)"
GORELEASER_LD_FLAGS = -s -w
GORELEASER_HOMEBREW_NAME = kubectl-falcon
GOLANG_CROSS_VERSION = v1.15.2
GORELEASER_CONFIG = .goreleaser.yml
GORELEASER_SKIP_VALIDATE ?= false

release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e MAINNET=$(MAINNET) \
		-e BUILD_FLAGS="$(GORELEASER_FLAGS)" \
		-e LD_FLAGS="$(GORELEASER_LD_FLAGS)" \
		-e HOMEBREW_NAME="$(GORELEASER_HOMEBREW_NAME)" \
		-e HOMEBREW_CUSTOM="$(GORELEASER_HOMEBREW_CUSTOM)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/github.com/crowdstrike/kubectl-falcon \
		-w /go/src/github.com/crowdstrike/kubectl-falcon \
		troian/golang-cross:${GOLANG_CROSS_VERSION} \
		-f "$(GORELEASER_CONFIG)" --skip-validate=$(GORELEASER_SKIP_VALIDATE) --rm-dist --skip-publish
