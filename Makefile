.PHONY: test test-cover test-cover-all test-cover-html

# Полное покрытие всего репозитория (включая main, docs, DI, HTTP handlers, SQL repos).
test-cover-all:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	@echo "--- total (all packages) ---"
	@go tool cover -func=coverage.out | tail -1

# Unit coverage: только пакеты с бизнес-логикой, которые покрываются unit-тестами.
# Исключены: cmd, docs, DI-контейнер, HTTP controllers, workers, cache, SQL repositories.
test-cover:
	go test \
		./internal/domain/... \
		./internal/application/... \
		./internal/errors/... \
		./internal/infrastructure/repository/common/... \
		./internal/endpoint/controller/http \
		./pkg/apperror/... \
		./pkg/response/... \
		./pkg/token/... \
		./pkg/util/... \
		./pkg/restclient/... \
		./pkg/cron/... \
		./pkg/database/postgres/... \
		./config/... \
		-coverprofile=coverage.out \
		-covermode=atomic
	@echo "--- total (unit scope) ---"
	@go tool cover -func=coverage.out | tail -1

test-cover-html: test-cover
	go tool cover -html=coverage.out -o coverage.html
	@echo "report: coverage.html"

test:
	go test ./...

test-integration:
	go test -tags=integration ./tests/integration/... -v -count=1 -timeout=10m

test-all: test test-integration
