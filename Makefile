.PHONY: lint
lint:
	buf lint
	api-linter
