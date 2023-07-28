.PHONY: tidy
tidy:
	for app in config api user auth gateway; do $(MAKE) -C $${app} tidy; done

.PHONY: gen
gen: gen-proto
	$(MAKE) -C user gen

.PHONY: gen-proto
gen-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=./api/ --go-grpc_out=./api/ --proto_path=./api user.proto auth.proto

.PHONY: build
build:
	for app in user auth gateway; do $(MAKE) -C $${app} build; done
