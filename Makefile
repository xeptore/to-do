.PHONY: tidy
tidy:
	for app in user auth; do $(MAKE) -C $${app} tidy; done

.PHONY: gen
gen:
	for app in user auth; do $(MAKE) -C $${app} gen; done

.PHONY: build
build:
	for app in user auth; do $(MAKE) -C $${app} build; done
