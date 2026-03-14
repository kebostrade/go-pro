# Protocol Buffer Compilation Makefile

.PHONY: all clean help

# Variables
PROTO_FILES := *.proto
GO_OUT := .

# Colors
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_BLUE := \033[34m

all: generate

## generate: Generate Go code from .proto files
generate:
	@echo "$(COLOR_BLUE)Generating Go code from .proto files...$(COLOR_RESET)"
	@protoc --go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "$(COLOR_GREEN)✓ Generated: $(PROTO_FILES)$(COLOR_RESET)"

## clean: Remove generated files
clean:
	@echo "$(COLOR_BLUE)Cleaning generated files...$(COLOR_RESET)"
	@rm -f *.pb.go
	@echo "$(COLOR_GREEN)✓ Cleaned generated files$(COLOR_RESET)"

## help: Show help message
help:
	@echo "$(COLOR_BLUE)Protocol Buffer Compilation$(COLOR_RESET)"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
