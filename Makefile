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

.PHONY: test
test:
	@echo "Starting unit tests"
	@go test -v -conver ./...

.PHONY: mock
mock:
	@echo "Generating mocks"
	@mockgen --build_flags=--mod=mod -destination=internal/account/mock_repository.go -package account github.com/mehmetokdemir/currency-conversion-service/internal/account IAccountRepository
	@mockgen --build_flags=--mod=mod -destination=internal/account/mock_service.go -package account github.com/mehmetokdemir/currency-conversion-service/internal/account IAccountService
	@mockgen --build_flags=--mod=mod -destination=internal/user/mock_repository.go -package user github.com/mehmetokdemir/currency-conversion-service/internal/user IUserRepository
	@mockgen --build_flags=--mod=mod -destination=internal/user/mock_service.go -package user github.com/mehmetokdemir/currency-conversion-service/internal/user IUserService
	@mockgen --build_flags=--mod=mod -destination=internal/exchange/mock_repository.go -package exchange github.com/mehmetokdemir/currency-conversion-service/internal/exchange IExchangeRepository
	@mockgen --build_flags=--mod=mod -destination=internal/exchange/mock_service.go -package exchange github.com/mehmetokdemir/currency-conversion-service/internal/exchange IExchangeService

.PHONY: build
build: tidy
	@echo "[BUILD] Building the service"
	go build -o bin/currency-conversion-service

.PHONY: run
run:
	@echo "[RUN] Running the service"
	./bin/currency-conversion-service

.PHONY: rwgocommand
rwgocommand:
	@echo "[RUN] Running the service"
	go run .

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