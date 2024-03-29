SHELL := bash

# Directory, where all required tools are located (absolute path required)
HACK_DIR    ?= $(shell cd hack 2>/dev/null && pwd)

VERSION   ?= 0.0.1-local
KUBE_NAMESPACE ?= platform-monoskope-monoskope

LATEST_REV=$(shell git rev-list --tags --max-count=1)
LATEST_TAG=$(shell git describe --tags $(LATEST_REV))

export 

# go
include *go.mk

add-license: ## Adds the license to every file
	@docker run --rm -v "$(PWD):/src" -u $(shell id -u) ghcr.io/google/addlicense -c "Monoskope Authors" -l "apache" -v .

check-license: ## Checks thath the license is set on every file
	@docker run --rm -v "$(PWD):/src" -u $(shell id -u) ghcr.io/google/addlicense -c "Monoskope Authors" -l "apache" -v -check .

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
