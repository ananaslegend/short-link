.PHONY: generate
generate:
	@echo "Generating protobuf code..."
	@protoc --go_out=gen --go_opt=paths=source_relative \
		-I . events/v1/*.proto
	@echo "Done."
