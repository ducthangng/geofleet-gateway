PROTO_DIR = proto/user
OUT_DIR = pb/user

consul:
	docker run -d --name consul-dev -p 8500:8500 -p 8600:8600/udp consul:1.15.4

gen-proto:
	@mkdir -p $(OUT_DIR)
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/*.proto
	
clean-proto:
	rm pb/*.go

update-proto:
	GOPROXY=direct go get github.com/ducthangng/geofleet-proto@latest
	go mod tidy


.PHONY: consul update-proto