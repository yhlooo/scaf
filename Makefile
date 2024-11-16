LOCALBIN ?= $(shell pwd)/bin

.PHONY: grpc-gen
grpc-gen: protoc-gen-go protoc-gen-go-grpc
	protoc \
      --plugin="$(LOCALBIN)/protoc-gen-go" \
      --plugin="$(LOCALBIN)/protoc-gen-go-grpc" \
      --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      pkg/apis/*/*/grpc/*.proto

# 安装 protoc-gen-go
.PHONY: protoc-gen-go
protoc-gen-go: $(LOCALBIN)/protoc-gen-go
$(LOCALBIN)/protoc-gen-go: $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1

# 安装 protoc-gen-go-grpc
.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc: $(LOCALBIN)/protoc-gen-go-grpc
$(LOCALBIN)/protoc-gen-go-grpc: $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"
