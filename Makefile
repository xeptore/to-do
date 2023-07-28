.PHONY: tidy
tidy:
	for app in todo config api user auth gateway; do $(MAKE) -C $${app} tidy; done

.PHONY: up
up:
	for app in todo config api user auth gateway; do $(MAKE) -C $${app} up; done

.PHONY: gen
gen: gen-proto
	$(MAKE) -C user gen

.PHONY: gen-proto
gen-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=./api/pb --go-grpc_out=./api/pb --proto_path=./api user.proto auth.proto todo.proto

.PHONY: build
build:
	for app in todo user auth gateway; do $(MAKE) -C $${app} build; done
