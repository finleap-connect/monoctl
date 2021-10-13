BUILD_PATH ?= $(shell pwd)
GO_MODULE_MONOSKOPE ?= github.com/finleap-connect/monoskope
GO_MODULE ?= gitlab.figo.systems/platform/monoskope/monoctl

GO             ?= go

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.14.2

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.36.0

MOCKGEN         ?= $(TOOLS_DIR)/mockgen
GOMOCK_VERSION  ?= v1.5.0

COMMIT     	   := $(shell git rev-parse --short HEAD)
LDFLAGS    	   += -X=$(GO_MODULE)/internal/version.Version=$(VERSION) -X=$(GO_MODULE)/internal/version.Commit=$(COMMIT)
BUILDFLAGS 	   += -installsuffix cgo --tags release
PROTOC     	   ?= protoc

CMD_MONOCTL_LINUX = $(BUILD_PATH)/monoctl-linux-amd64
CMD_MONOCTL_OSX = $(BUILD_PATH)/monoctl-osx-amd64
CMD_MONOCTL_WIN = $(BUILD_PATH)/monoctl-win-amd64
CMD_MONOCTL_SRC = cmd/monoctl/*.go

OS = windows
ifeq ($(uname_S),Linux)
OS = linux
endif
ifeq ($(uname_S),Darwin)
OS = darwin
endif

.PHONY: go lint mod vet test clean build-monoctl-linux build-monoctl-all push-monoctl protobuf

mod: ## Do go mod tidy, download, verify
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod verify

vet: ## Do go ver
	$(GO) vet ./...

lint: ## Do golangci-lint
	$(LINTER) run -v --no-config --deadline=5m

go: mod vet lint test ## Do go mod / vet / lint /test

run: ## run monoctl, use `ARGS="get user"` to pass arguments
	$(GO) run -ldflags "$(LDFLAGS)" cmd/monoctl/*.go $(ARGS)

test: ## run all tests
	@find . -name '*.coverprofile' -exec rm {} \;
	$(GINKGO) -r -v -cover *
	@echo "mode: set" > ./monoctl.coverprofile
	@find ./internal -name "*.coverprofile" -exec cat {} \; | grep -v mode: | sort -r >> ./monoctl.coverprofile   
	@find ./internal -name '*.coverprofile' -exec rm {} \;

coverage: ## show test coverage
	@find . -name '*.coverprofile' -exec go tool cover -func {} \;

loc: ## show loc statistics
	@gocloc .

ginkgo-get: ## download ginkgo
	$(shell $(TOOLS_DIR)/goget-wrapper github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get: ## download golangci-lint
	$(shell curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) $(LINTER_VERSION))

gomock-get: ## download gomock
	$(shell $(TOOLS_DIR)/goget-wrapper github.com/golang/mock/mockgen@$(GOMOCK_VERSION))

ginkgo-clean: ## cleanup ginkgo
	rm -Rf $(TOOLS_DIR)/ginkgo

golangci-lint-clean: ## cleanup golangci-lint
	rm -Rf $(TOOLS_DIR)/golangci-lint

tools: golangci-lint-get ginkgo-get gomock-get  ## Target to install all required tools into TOOLS_DIR

clean: ginkgo-clean golangci-lint-clean build-clean  ## Target clean up tools in TOOLS_DIR
	rm -Rf reports/
	find . -name '*.coverprofile' -exec rm {} \;

build-clean: ## clean up binaries
	rm -Rf $(CMD_MONOCTL_LINUX)
	rm -Rf $(CMD_MONOCTL_OSX)
	rm -Rf $(CMD_MONOCTL_WIN)

$(CMD_MONOCTL_LINUX): ## build monoctl for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_LINUX) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_OSX): ## build monoctl for osx
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_OSX) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_WIN): ## build monoctl for windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_WIN) $(CMD_MONOCTL_SRC)

build-monoctl-linux: $(CMD_MONOCTL_LINUX) ## build monoctl for linux

build-monoctl-all: $(CMD_MONOCTL_LINUX) $(CMD_MONOCTL_OSX) $(CMD_MONOCTL_WIN) ## build monoctl for linux, osx and windows

push-monoctl:  ## push monoctl to artifactory
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_LINUX) "https://artifactory.figo.systems/artifactory/binaries/linux/monoctl-$(VERSION)"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_LINUX) "https://artifactory.figo.systems/artifactory/binaries/linux/monoctl"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_OSX) "https://artifactory.figo.systems/artifactory/binaries/osx/monoctl-$(VERSION)"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_OSX) "https://artifactory.figo.systems/artifactory/binaries/osx/monoctl"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_WIN) "https://artifactory.figo.systems/artifactory/binaries/win/monoctl-$(VERSION)"
	@curl -u$(ARTIFACTORY_BINARY_USER):$(ARTIFACTORY_BINARY_PW) -T $(CMD_MONOCTL_WIN) "https://artifactory.figo.systems/artifactory/binaries/win/monoctl"

rebuild-mocks: ## rebuild go mocks
	$(MOCKGEN) -package eventsourcing -destination test/mock/eventsourcing/command_handler_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing CommandHandlerClient
	$(MOCKGEN) -package domain -destination test/mock/domain/cluster_client.go github.com/finleap-connect/monoskope/pkg/api/domain ClusterClient,Cluster_GetAllClient
	$(MOCKGEN) -package domain -destination test/mock/domain/tenant_client.go github.com/finleap-connect/monoskope/pkg/api/domain TenantClient,Tenant_GetAllClient
	$(MOCKGEN) -package domain -destination test/mock/domain/certificate_client.go github.com/finleap-connect/monoskope/pkg/api/domain CertificateClient
	$(MOCKGEN) -package domain -destination test/mock/gateway/cluster_auth_client.go github.com/finleap-connect/monoskope/pkg/api/gateway ClusterAuthClient
