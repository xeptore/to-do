.PHONY: gen-proto
gen-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=./internal --go-grpc_out=./internal --proto_path=../api user.proto

.PHONY: gen
gen: gen-proto gen-db

.PHONY: clean
clean:
	go clean -r -cache -testcache -modcache

.PHONY: tidy
tidy:
	go mod tidy -v -x

.PHONY: gen-db
gen-db:
	go run github.com/joho/godotenv/cmd/godotenv@latest -f .env go generate ./db/gen/gen.go

.PHONY: build-clean
build-clean: clean build

.PHONY: test
test:
	go test -trimpath -buildvcs=false -ldflags '-extldflags "-static" -s -w -buildid=' -race -failfast -vet=all -covermode=atomic -coverprofile=coverage.out -v ./...

ifndef app_version
app_version := dev
endif
.PHONY: build
build:
	rm -rf ./bin
	mkdir -p ./bin
	go build -trimpath -buildvcs=false -ldflags "-extldflags '-static' -s -w -buildid=''" -o ./bin/user .
