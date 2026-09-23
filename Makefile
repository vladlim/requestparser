.PHONY: check fmt fmt-check vet test

check: fmt-check vet test

fmt:
	gofmt -w .

fmt-check:
	@files=$$(gofmt -l .) || exit 1; \
	if [ -n "$$files" ]; then \
		printf 'Run make fmt on these files:\n%s\n' "$$files"; \
		exit 1; \
	fi

vet:
	go vet ./...

test:
	go test ./...
