.PHONY: tidy
tidy:
	go mod tidy -v -x

.PHONY: up
up:
	go get -u -d -d ./...
	$(MAKE) tidy
