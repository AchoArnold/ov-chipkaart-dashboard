.PHONY: compile
compile: ## Compile the proto file.
	protoc -I shared/proto/transactions/ shared/proto/transactions/raw_records.proto --go_out=plugins=grpc:shared/proto/transactions/