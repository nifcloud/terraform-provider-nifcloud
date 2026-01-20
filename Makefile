ifneq ($(shell command -v git 2> /dev/null),)
TAG          ?= $(shell git describe --tags --abbrev=0 2>/dev/null)
ON_EXACT_TAG ?= $(shell git name-rev --name-only --tags --no-undefined HEAD 2>/dev/null | sed -n 's/^\([^^~]\{1,\}\)\(\^0\)\{0,1\}$$/\1/p')
VERSION      ?= $(shell [ -z "$(ON_EXACT_TAG)" ] && echo "v0.1.0" | sed 's/^v//' || echo "$(TAG)" | sed 's/^v//' )
else
VERSION ?= "0.1.0"
endif

PROVIDER=nifcloud
GOBIN=$(shell pwd)/bin
BINARY=$(GOBIN)/terraform-provider-$(PROVIDER)
PLUGIN_0.13=registry.terraform.io/nifcloud/nifcloud/$(VERSION)/$(shell uname -s|tr '[:upper:]' '[:lower:]')_amd64/terraform-provider-nifcloud_v$(VERSION)

default: build

##################
# Install tools  #
##################
.PHONY: tools
tools:
	@ mkdir -p ${GOBIN}
	@ GOBIN=${GOBIN} go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint
	@ GOBIN=${GOBIN} go install github.com/bflad/tfproviderdocs
	@ GOBIN=${GOBIN} go install github.com/client9/misspell/cmd/misspell
	@ GOBIN=${GOBIN} go install github.com/katbyte/terrafmt
	@ GOBIN=${GOBIN} go install github.com/bflad/tfproviderlint/cmd/tfproviderlint

####################
# Building/Install #
####################
.PHONY: build
build:
	@ mkdir -p ${GOBIN}
	@ go build -o ${BINARY}
.PHONY: install
install: build
	@ mkdir -p ~/.terraform.d/plugins
	@ cp ${BINARY} ~/.terraform.d/plugins
	@ mkdir -p $(shell dirname ~/.terraform.d/plugins/$(PLUGIN_0.13))
	@ cp ${BINARY} ~/.terraform.d/plugins/$(PLUGIN_0.13)

###########
# Testing #
###########
.PHONY: test
test:
	@ go test ./$(PROVIDER)/... -timeout 120s -p 4 -race -cover -coverprofile=coverage.out
.PHONY: testacc
testacc:
	@ TF_ACC=1 go test ./$(PROVIDER)/acc/... -v -count 1 -parallel 4 -timeout 360m -run Test -coverprofile=coverage.out -coverpkg=./...
.PHONY: sweep
sweep:
	@ echo "-> WARNING: This will destroy infrastructure. Use only in development accounts."
	@ read -r -p "do you wish to continue? [y/N]: " res && if [[ "$${res:0:1}" =~ ^([yY]) ]]; then echo "-> Continuing..."; else exit 1; fi
	@ go test ./$(PROVIDER)/acc/... -v -sweep=$(NIFCLOUD_DEFAULT_REGION) -timeout 60m

##################
# Linting/Verify #
##################
.PHONY: lint
lint:
	@ ${GOBIN}/golangci-lint run ./$(PROVIDER)/...
	@ $(GOBIN)/tfproviderlint ./$(PROVIDER)/...
.PHONY: docscheck
docscheck:
	@ find ./docs -type f -name "*.md" -exec $(GOBIN)/terrafmt diff -c -q {} \;
	@ $(GOBIN)/tfproviderdocs check -provider-name $(PROVIDER)
	@ $(GOBIN)/misspell -error -source=text ./docs/
validate-examples:
	@ for example in $(shell find examples -maxdepth 1 -type d | grep "/"); do \
		cd $(shell pwd)/$$example ; \
		terraform init ; \
		terraform validate ; \
	  done
