.PHONY: tidy
tidy:
	for app in config user auth gateway; do $(MAKE) -C $${app} tidy; done

.PHONY: gen
gen:
	for app in user auth gateway; do $(MAKE) -C $${app} gen; done

.PHONY: build
build:
	for app in user auth gateway; do $(MAKE) -C $${app} build; done
