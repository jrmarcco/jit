# initialize project environment
.PHONY: setup
setup:
	@sh ./scripts/setup_env.sh
	@sh ./scripts/setup_hook.sh

# format code
.PHONY: fmt
fmt:
	@goimports -l -w $$(find . -type f -name '*.go' -not -path "./.idea/*" -not -path "./.vscode/*")
	@gofumpt -l -w $$(find . -type f -name '*.go' -not -path "./.idea/*" -not -path "./.vscode/*")

# clean up project dependencies
.PHONY: tidy
tidy:
	@go mod tidy -v

# check code format and dependencies
.PHONY: check
check:
	@$(MAKE) --no-print-directory fmt
	@$(MAKE) --no-print-directory tidy

# static check with go vet
.PHONY: vet
vet:
	@go vet ./...

# lint code (baseline tier, CI blocking)
.PHONY: lint
lint:
	@golangci-lint run -c ./scripts/lint/.golangci.yaml ./...

# lint code (strict tier, report-only target)
.PHONY: lint-strict
lint-strict:
	@golangci-lint run -c ./scripts/lint/.golangci.strict.yaml ./...

# list tags, newest first
.PHONY: tags
tags:
	@git tag --sort=-v:refname | head -20

# tag the current commit and push it (fires the release pipeline):
#   make tag VERSION=1.0.6           -> annotated v1.0.6 with a default message
#   make tag VERSION=1.0.6 MSG="..." -> with your own message
# VERSION keeps its bare form: a leading v is accepted and stripped.
.PHONY: tag
tag:
	@test -n "$(VERSION)" || { echo "usage: make tag VERSION=<x.y.z> [MSG=...]"; exit 2; }
	@v=$$(printf '%s' "$(VERSION)" | sed 's/^v//'); \
	if [ -n "$$(git tag -l v$$v)" ]; then echo "v$$v already exists"; exit 1; fi; \
	git tag -a "v$$v" -m "$(if $(MSG),$(MSG),v$$v)" && \
	git push origin "v$$v"

# bump from the newest tag and push it (no VERSION to type):
#   make tag-patch [MSG=...]   v1.0.5 -> v1.0.6
#   make tag-minor [MSG=...]   v1.0.5 -> v1.1.0
#   make tag-major [MSG=...]   v1.0.5 -> v2.0.0
# pre-release tags (-rc1 and friends) are skipped when picking the newest.
.PHONY: tag-patch tag-minor tag-major
tag-patch tag-minor tag-major:
	@t="$@"; set -- $$(git tag --sort=-v:refname | sed -n 's/^v\([0-9]\+\)\.\([0-9]\+\)\.\([0-9]\+\)$$/\1 \2 \3/p' | head -1); \
	if [ -z "$$1" ]; then echo "no x.y.z tag to bump"; exit 1; fi; \
	case "$$t" in \
		tag-patch) set -- "$$1" "$$2" "$$(($$3 + 1))";; \
		tag-minor) set -- "$$1" "$$(($$2 + 1))" 0;; \
		tag-major) set -- "$$(($$1 + 1))" 0 0;; \
	esac; \
	v="$$1.$$2.$$3"; \
	echo "v$$v it is"; \
	if [ -n "$$(git tag -l v$$v)" ]; then echo "v$$v already exists"; exit 1; fi; \
	git tag -a "v$$v" -m "$(if $(MSG),$(MSG),v$$v)" && \
	git push origin "v$$v"

# drop a tag everywhere (local and on origin) — for a misfire before the release runs:
#   make untag VERSION=1.0.6
.PHONY: untag
untag:
	@test -n "$(VERSION)" || { echo "usage: make untag VERSION=<x.y.z>"; exit 2; }
	@v=$$(printf '%s' "$(VERSION)" | sed 's/^v//'); \
	git tag -d "v$$v" 2>/dev/null || true; \
	{ git push origin ":v$$v" 2>/dev/null || true; }; \
	echo "v$$v dropped (if it was anywhere)"
