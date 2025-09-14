# Directory containing .proto files
PROTO_DIR = proto

# Collect all proto files (works cross-platform if filenames are simple)
PROTO_SRC = $(wildcard $(PROTO_DIR)/*.proto)

# Output directory for generated Go code
GO_OUT = .

.PHONY: generate-proto clean

# Generate Go and gRPC code from .proto files
generate-proto:
	@echo "Generating gRPC code..."
	protoc --proto_path=$(PROTO_DIR) --go_out=$(GO_OUT) --go-grpc_out=$(GO_OUT) $(PROTO_SRC)

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	-del /Q $(GO_OUT)\*.pb.go 2> NUL || rm -f $(GO_OUT)/*.pb.go
