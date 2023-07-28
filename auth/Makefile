.PHONY: clean
clean:
	go clean -r -cache -testcache -modcache

.PHONY: tidy
tidy:
	go mod tidy -v -x

.PHONY: up
up:
	go get -u -d -d ./...
	$(MAKE) tidy

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
	go build -trimpath -buildvcs=false -ldflags "-extldflags '-static' -s -w -buildid=''" -o ./bin/auth .
