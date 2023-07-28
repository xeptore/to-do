.PHONY: tidy
tidy:
	go mod tidy -v -x

.PHONY: up
up:
	go list  -f '{{ range .Imports }}{{ printf "%s\n" . }}{{ end }}' ./... | grep github.com/xeptore/to-do/ | xargs -r go get -u -d -d
	$(MAKE) tidy
