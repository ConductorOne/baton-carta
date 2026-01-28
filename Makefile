GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
BUILD_DIR = dist/${GOOS}_${GOARCH}
GENERATED_CONF := pkg/config/conf.gen.go
OUTPUT_PATH = ${BUILD_DIR}/baton-carta

# Set the build tag conditionally based on BATON_LAMBDA_SUPPORT
ifdef BATON_LAMBDA_SUPPORT
	BUILD_TAGS=-tags baton_lambda_support
else
	BUILD_TAGS=
endif

$(GENERATED_CONF): pkg/config/config.go go.mod
	@echo "Generating $(GENERATED_CONF)..."
	go generate ./pkg/config

.PHONY: generate
generate: $(GENERATED_CONF)

.PHONY: build
build: $(GENERATED_CONF)
	rm -f ${OUTPUT_PATH}
	mkdir -p ${BUILD_DIR}
	go build ${BUILD_TAGS} -o ${OUTPUT_PATH} cmd/baton-carta/*.go

.PHONY: update-deps
update-deps:
	go get -d -u ./...
	go mod tidy -v
	go mod vendor

.PHONY: add-dep
add-dep:
	go mod tidy -v
	go mod vendor

.PHONY: lint
lint:
	golangci-lint run
