BUILD_PATH ?= $(shell pwd)
GO_MODULE_MONOSKOPE ?= github.com/finleap-connect/monoskope
GO_MODULE ?= github.com/finleap-connect/monoctl

GO             ?= go
GOGET          ?= $(HACK_DIR)/goget-wrapper

GINKGO         ?= $(TOOLS_DIR)/ginkgo
GINKO_VERSION  ?= v1.16.4

LINTER 	   	   ?= $(TOOLS_DIR)/golangci-lint
LINTER_VERSION ?= v1.39.0

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

.PHONY: go lint mod vet test clean build-monoctl-linux build-monoctl-all protobuf

mod: ## Do go mod tidy, download, verify
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod verify

vet: ## Do go ver
	$(GO) vet ./...

go: mod vet lint test ## Do go mod / vet / lint /test

run: ## run monoctl, use `ARGS="get user"` to pass arguments
	$(GO) run -ldflags "$(LDFLAGS)" cmd/monoctl/*.go $(ARGS)

test: ## run all tests
# https://onsi.github.io/ginkgo/#running-tests
	find . -name '*.coverprofile' -exec rm {} \;
	@$(GINKGO) -r -v -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoctl.coverprofile 

test-ci: ## run all tests in CICD
# https://onsi.github.io/ginkgo/#running-tests
	find . -name '*.coverprofile' -exec rm {} \;
	@$(GINKGO) -r -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoctl.coverprofile 

coverage: ## print coverage from coverprofiles
	@go tool cover -func monoctl.coverprofile

ginkgo-get $(GINKGO):
	$(shell $(GOGET) github.com/onsi/ginkgo/ginkgo@$(GINKO_VERSION))

golangci-lint-get $(LINTER):
	$(shell $(HACK_DIR)/golangci-lint.sh -b $(TOOLS_DIR) $(LINTER_VERSION))

gomock-get: ## download gomock
	$(shell $(TOOLS_DIR)/goget-wrapper github.com/golang/mock/mockgen@$(GOMOCK_VERSION))

lint: $(LINTER) ## go lint
	$(LINTER) run -v -E goconst -E misspell

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

rebuild-mocks: ## rebuild go mocks
	$(MOCKGEN) -package eventsourcing -destination test/mock/eventsourcing/command_handler_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing CommandHandlerClient
	$(MOCKGEN) -package domain -destination test/mock/domain/cluster_client.go github.com/finleap-connect/monoskope/pkg/api/domain ClusterClient,Cluster_GetAllClient
	$(MOCKGEN) -package domain -destination test/mock/domain/tenant_client.go github.com/finleap-connect/monoskope/pkg/api/domain TenantClient,Tenant_GetAllClient
	$(MOCKGEN) -package domain -destination test/mock/domain/certificate_client.go github.com/finleap-connect/monoskope/pkg/api/domain CertificateClient
	$(MOCKGEN) -package domain -destination test/mock/gateway/cluster_auth_client.go github.com/finleap-connect/monoskope/pkg/api/gateway ClusterAuthClient
	$(MOCKGEN) -package domain -destination test/mock/gateway/api_token_client.go github.com/finleap-connect/monoskope/pkg/api/gateway APITokenClient
	$(MOCKGEN) -package domain -destination test/mock/domain/audit_log_client.go github.com/finleap-connect/monoskope/pkg/api/domain AuditLogClient,AuditLog_GetByDateRangeClient,AuditLog_GetUserActionsClient