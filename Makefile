.DEFAULT_GOAL := help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

OS := $(shell uname)
MODULE := github.com/TrogonStack/TrogonEventStore-Client-Go
PROTO_ROOT := protos
PROTO_FILES := $(shell find $(PROTO_ROOT)/trogoneventstore -type f -name '*.proto' | sort)
# google.rpc bindings are consumed from google.golang.org/genproto, as declared by their go_package options.
EXTERNAL_GO_PROTO_FILES := $(PROTO_ROOT)/trogoneventstore/protocols/v1/code.proto $(PROTO_ROOT)/trogoneventstore/protocols/v1/status.proto
CLIENT_PROTO_FILES := $(filter-out $(EXTERNAL_GO_PROTO_FILES),$(PROTO_FILES))
PROTOC_GEN_GO_VERSION := v1.36.9
PROTOC_GEN_GO_GRPC_VERSION := v1.5.1

.PHONY: build
build: ## Compile all packages.
	go build ./...

.PHONY: validate-protos
validate-protos: ## Compile every checked-in protobuf source.
	@descriptor=$$(mktemp); \
	trap 'rm -f "$$descriptor"' EXIT; \
	protoc -I $(PROTO_ROOT) --include_imports --descriptor_set_out="$$descriptor" $(PROTO_FILES)

.PHONY: generate-protos
generate-protos: validate-protos ## Regenerate repository-owned protobuf and gRPC bindings.
	@tool_dir=$$(mktemp -d); \
	trap 'rm -rf "$$tool_dir"' EXIT; \
	GOBIN="$$tool_dir" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION); \
	GOBIN="$$tool_dir" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION); \
	PATH="$$tool_dir:$$PATH" protoc -I $(PROTO_ROOT) \
		--go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		$(CLIENT_PROTO_FILES)

.PHONY: verify-protos
verify-protos: generate-protos ## Fail when generated protobuf bindings drift.
	@git diff --exit-code -- $(PROTO_ROOT)

.PHONY: start-server
start-server: ## Start the TrogonEventStore integration-test cluster.
	@docker --version
	@docker compose up -d
	@for endpoint in http://localhost:2114/-/liveness https://localhost:2115/-/liveness; do \
		attempt=0; \
		until curl --insecure --fail --silent "$$endpoint" >/dev/null; do \
			attempt=$$((attempt + 1)); \
			[ "$$attempt" -lt 60 ] || exit 1; \
			sleep 1; \
		done; \
	done

.PHONY: stop-server
stop-server: ## Stop the TrogonEventStore integration-test cluster.
	@docker compose down -v --remove-orphans

.PHONY: start-oauth
start-oauth: ## Start the OAuth integration-test stack.
	@docker compose -f docker-compose.oauth.yml up -d --wait

.PHONY: stop-oauth
stop-oauth: ## Stop the OAuth integration-test stack.
	@docker compose -f docker-compose.oauth.yml down -v --remove-orphans

.PHONY: test
test: ## Run all tests.
	go test --count=1 ./...
