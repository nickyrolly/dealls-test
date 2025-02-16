# include all prerequisites
include scripts/make/*.Mak

.bin-validator:
	go build -o app-dealls-test cmd/app/*.go  # Adjust the path

init: .bin-validator
	# Other initialization commands, if any
	@touch dealls_test.db
	@go mod tidy
	@sleep 5
	@make migrate

.PHONY: test
test:
	go test -v -race -cover ./...