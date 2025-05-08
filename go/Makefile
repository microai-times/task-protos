# Makefile for Protocol Buffers compilation

# Variables
PROTO_DIR := ../protos
OUT_DIR := gen
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
PROTO_TARGETS := $(patsubst $(PROTO_DIR)/%.proto,$(OUT_DIR)/%.pb.go,$(PROTO_FILES))

# Default target
all: $(PROTO_TARGETS)

# Protobuf compilation rule
$(OUT_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(OUT_DIR)
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$^

# Clean generated files
clean:
	rm -rf $(OUT_DIR)

.PHONY: all clean