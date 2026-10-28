.PHONY: build test check
build:
	CGO_ENABLED=0 go build -o bin/factoryctl ./cmd/factoryctl
test:
	go test ./...
check:
	go test -race ./...
	go vet ./...

# Accounting and full E2E qualification are intentionally separate gates.
REVIEW_EVIDENCE ?= artifacts/cli-review/848b4ec/requirements-coverage.json
.PHONY: review-check review-gate
review-check:
	node scripts/review-video/check-matrix.mjs $(REVIEW_EVIDENCE)
review-gate:
	node scripts/review-video/check-matrix.mjs $(REVIEW_EVIDENCE) --qualify

.PHONY: install-hooks repository-check
install-hooks:
	git config core.hooksPath .githooks
repository-check:
	python3 scripts/check-repository-content.py --staged
	python3 -m unittest discover -s scripts -p test_repository_content.py
