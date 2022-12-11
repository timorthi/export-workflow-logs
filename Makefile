.PHONY: test
test:
	@echo "Test placeholder"

.PHONY: build
build:
	docker build -t export-workflow-logs .
