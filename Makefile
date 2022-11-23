include META

#LD flags
LDFLAGS := -w

#Environment Settings
BUILD_ENV_SET_FOR_LINUX := $(RUN_ENV_SET) GOOS=linux


.PHONY: generate-docs
generate-docs:
	@echo "[GENERATE DOCS] Generating API documents"
	@echo " - Updating document version"
	@echo " - Initializing swag"
	@swag init --parseDependency --parseInternal --generatedTime --parseDepth 3

.PHONY: tidy
tidy:
	@echo "[TIDY] Running go mod tidy"
	@$(RUN_ENV_SET) go mod tidy -compat=1.19

.PHONY: lint
lint:
	@echo "[TIDY] Running golangci-lint run"
	@golangci-lint run


.PHONY: build
build: tidy
	@echo "[BUILD] Building the service"
	@$(BUILD_ENV_SET_FOR_LINUX) go build -mod=mod -installsuffix cgo -ldflags '$(LDFLAGS)' -o bin/service .

.PHONY: git
git:
	@echo "[BUILD] Committing and pushing to remote repository"
	@echo " - Committing"
	@git add META
	@git commit -am "v$(VERSION)"
	@echo " - Tagging"
	@git tag v${VERSION}
	@echo " - Pushing"
	@git push --tags origin ${BRANCH}