.PHONY: test
test:
	go test -v

.PHONY: build
build:
	docker build -t export-workflow-logs .
