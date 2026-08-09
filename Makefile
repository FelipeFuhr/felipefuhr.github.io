# felipefuhr.github.io — The workshop behind ffreis.com (static site, built by a small Go generator).
SHELL := /usr/bin/env bash
SHELL_SCRIPTS := $(shell find scripts/ -name '*.sh' ! -type l 2>/dev/null | sort)
LEFTHOOK_VERSION ?= 1.7.10
LEFTHOOK_BIN     ?= $(CURDIR)/.bin/lefthook
QUALITY_KIT_SCRIPTS ?= /media/ffreis/second/projects/quality-kit/scripts
COVERAGE_MIN ?= 75
MUTATION_PACKAGES  ?= ./internal/site/...
MUTATION_THRESHOLD ?= 60
GITLEAKS ?= gitleaks

.PHONY: help ci build serve fmt fmt-check lint vet shellcheck test coverage-gate mutation clean ci-local init-github lefthook-bootstrap secrets-scan-staged

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

ci: fmt-check lint coverage-gate build ## Full local gate (matches CI)

build: ## Render the static site into dist/
	go run ./cmd/build

serve: build ## Build then serve dist/ at http://localhost:8080
	@echo "serving http://localhost:8080 (Ctrl-C to stop)"
	@cd dist && python3 -m http.server 8080

fmt: ## Format Go sources
	gofmt -w .

fmt-check: ## Fail if any Go file needs formatting
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

lint: vet shellcheck ## go vet + shellcheck

vet: ## go vet
	go vet ./...

shellcheck: ## shellcheck scripts/*.sh
	@command -v shellcheck >/dev/null 2>&1 && shellcheck -x $(SHELL_SCRIPTS) || echo "shellcheck not installed — skip"

test: ## go test (race + shuffle) with coverage profile
	go test -race -shuffle=on -coverprofile=coverage.out ./...

coverage-gate: test ## Enforce the coverage floor
	@bash scripts/hooks/check_coverage_gate.sh coverage.out $(COVERAGE_MIN)

mutation: ## Run mutation testing with gremlins on the tested core (slow — CI only)
	@which gremlins >/dev/null 2>&1 || go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
	gremlins unleash --threshold-efficacy $(MUTATION_THRESHOLD) $(MUTATION_PACKAGES)

# secrets-scan-staged was missing entirely: the shared lefthook `secret-scan`
# step (pinned to ffreis-platform-standards v1.6.0, which has no gitleaks
# fallback) runs `make secrets-scan-staged` unconditionally and failed outright
# on every commit. Every other fleet repo has this target — added to match
# convention (see e.g. ffreis-website-inventory/Makefile).
secrets-scan-staged: ## Scan staged diff for secrets with gitleaks
	@command -v $(GITLEAKS) >/dev/null 2>&1 || { echo "Missing required tool: $(GITLEAKS) — https://github.com/gitleaks/gitleaks#installing" >&2; exit 1; }
	$(GITLEAKS) protect --staged --redact

clean: ## Remove build + coverage artifacts
	rm -rf dist coverage.out

ci-local: ## Run workflows locally via the ci-local harness (GH Actions quota fallback). ARGS=...
	@mkdir -p scripts
	@curl -fsSL "https://raw.githubusercontent.com/FelipeFuhr/ffreis-platform-ci-local/v1.0.0/scripts/run-ci-local.sh" \
		-o scripts/run-ci-local.sh && chmod +x scripts/run-ci-local.sh
	@CI_LOCAL_FINDINGS_REF=v1.0.0 PATH="$(CURDIR)/.bin:$(PATH)" bash ./scripts/run-ci-local.sh $(ARGS)

init-github: ## Apply standard fleet settings to the GitHub repo
	@repo=$$(git remote get-url origin | sed -E 's|.*github\.com[:/]||; s|\.git$$||'); \
	bash "$(QUALITY_KIT_SCRIPTS)/configure-repo-settings.sh" "$$repo"

lefthook-bootstrap: ## Download the pinned lefthook binary + install hooks
	@bash scripts/bootstrap_lefthook.sh
