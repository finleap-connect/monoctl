BUILD_PATH ?= $(shell pwd)
GO_MODULE ?= github.com/finleap-connect/monoctl
GO_MODULE_MONOSKOPE ?= github.com/finleap-connect/monoskope

GO             ?= go

COMMIT     	   := $(shell git rev-parse --short HEAD)
LDFLAGS    	   += -X=$(GO_MODULE)/internal/version.Version=$(VERSION) -X=$(GO_MODULE)/internal/version.Commit=$(COMMIT)
BUILDFLAGS 	   += -installsuffix cgo --tags release

CMD_MONOCTL_LINUX = $(BUILD_PATH)/monoctl-linux-amd64
CMD_MONOCTL_LINUX_ARM = $(BUILD_PATH)/monoctl-linux-arm64
CMD_MONOCTL_OSX = $(BUILD_PATH)/monoctl-osx-amd64
CMD_MONOCTL_OSX_ARM = $(BUILD_PATH)/monoctl-osx-arm64
CMD_MONOCTL_WIN = $(BUILD_PATH)/monoctl-win-amd64
CMD_MONOCTL_SRC = cmd/monoctl/*.go

OS = windows
ifeq ($(uname_S),Linux)
OS = linux
endif
ifeq ($(uname_S),Darwin)
OS = darwin
endif

.PHONY: mod
mod: ## Do go mod tidy, download, verify
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod verify

.PHONY: fmt
fmt: ## Do go vet
	$(GO) fmt ./...

.PHONY: vet
vet: ## Do go vet
	$(GO) vet ./...

.PHONY: go
go: mod fmt vet lint test ## Do go mod / fmt / vet / lint /test

.PHONY: run
run: ## run monoctl, use `ARGS="get user"` to pass arguments
	$(GO) run -ldflags "$(LDFLAGS)" cmd/monoctl/*.go $(ARGS)

.PHONY: test
test: ginkgo ## run all tests
# https://onsi.github.io/ginkgo/#running-tests
	find . -name '*.coverprofile' -exec rm {} \;
	@$(GINKGO) -r -v -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoctl.coverprofile 

.PHONY: test-ci
test-ci: ## run all tests in CICD
# https://onsi.github.io/ginkgo/#running-tests
	find . -name '*.coverprofile' -exec rm {} \;
	@$(GINKGO) -r -cover --failFast -requireSuite -covermode count -outputdir=$(BUILD_PATH) -coverprofile=monoctl.coverprofile 

.PHONY: coverage
coverage: ## print coverage from coverprofiles
	@go tool cover -func monoctl.coverprofile

.PHONY: lint
lint: golangcilint ## go lint
	$(GOLANGCILINT) run -v -E goconst -E misspell -E gofmt

$(CMD_MONOCTL_LINUX): ## build monoctl for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_LINUX) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_LINUX_ARM): ## build monoctl for linux arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_LINUX_ARM) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_OSX): ## build monoctl for osx
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_OSX) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_OSX_ARM): ## build monoctl for osx arm
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_OSX_ARM) $(CMD_MONOCTL_SRC)

$(CMD_MONOCTL_WIN): ## build monoctl for windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a $(BUILDFLAGS) -ldflags "$(LDFLAGS) -X=$(GO_MODULE_MONOSKOPE)/pkg/logger.logMode=noop" -o $(CMD_MONOCTL_WIN) $(CMD_MONOCTL_SRC)

build-monoctl-linux: $(CMD_MONOCTL_LINUX) $(CMD_MONOCTL_LINUX_ARM) ## build monoctl for linux
	@chmod a+x $(CMD_MONOCTL_LINUX)
	@chmod a+x $(CMD_MONOCTL_LINUX_ARM)

build-monoctl-osx: $(CMD_MONOCTL_OSX) $(CMD_MONOCTL_OSX_ARM) ## build monoctl for osx
	@chmod a+x $(CMD_MONOCTL_OSX)
	@chmod a+x $(CMD_MONOCTL_OSX_ARM)

build-monoctl-win: $(CMD_MONOCTL_WIN)  ## build monoctl for windows
	@chmod a+x $(CMD_MONOCTL_WIN)

build-monoctl-all: build-monoctl-linux build-monoctl-osx build-monoctl-win ## build monoctl for linux, osx and windows

.PHONY: rebuild-mocks
rebuild-mocks: gomock ## rebuild go mocks
	$(MOCKGEN) -package eventsourcing -destination test/mock/eventsourcing/command_handler_client.go github.com/finleap-connect/monoskope/pkg/api/eventsourcing CommandHandlerClient
	$(MOCKGEN) -package domain -destination test/mock/domain/cluster_client.go github.com/finleap-connect/monoskope/pkg/api/domain ClusterClient,Cluster_GetAllClient,ClusterAccessClient,ClusterAccess_GetClusterAccessClient
	$(MOCKGEN) -package domain -destination test/mock/domain/tenant_client.go github.com/finleap-connect/monoskope/pkg/api/domain TenantClient,Tenant_GetAllClient
	$(MOCKGEN) -package domain -destination test/mock/gateway/cluster_auth_client.go github.com/finleap-connect/monoskope/pkg/api/gateway ClusterAuthClient
	$(MOCKGEN) -package domain -destination test/mock/gateway/api_token_client.go github.com/finleap-connect/monoskope/pkg/api/gateway APITokenClient
	$(MOCKGEN) -package domain -destination test/mock/domain/audit_log_client.go github.com/finleap-connect/monoskope/pkg/api/domain AuditLogClient,AuditLog_GetByDateRangeClient,AuditLog_GetByUserClient,AuditLog_GetUserActionsClient,AuditLog_GetUsersOverviewClient

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: clean
clean: ## Clean up build dependencies
	rm -R $(LOCALBIN)

## Tool Binaries
MOCKGEN ?= $(LOCALBIN)/mockgen
GINKGO ?= $(LOCALBIN)/ginkgo
GOLANGCILINT ?= $(LOCALBIN)/golangci-lint

## Tool Versions
GOMOCK_VERSION  ?= v1.5.0
GINKGO_VERSION ?= v1.16.5
GOLANGCILINT_VERSION ?= v1.48.0

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/ginkgo@$(GINKGO_VERSION)

.PHONY: golangcilint
golangcilint: $(GOLANGCILINT) ## Download golangci-lint locally if necessary.
$(GOLANGCILINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

.PHONY: gomock
gomock: $(MOCKGEN) ## Download mockgen locally if necessary.
$(MOCKGEN): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golang/mock/mockgen@$(GOMOCK_VERSION)
