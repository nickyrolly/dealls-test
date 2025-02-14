################################################################################
##@ PLAYGROUND
################################################################################

run: .bin-validator ## Run service development
	@export ENV=development && go run cmd/app/*.go

migrate: .bin-validator ## Migrate database "make migrate env=beta"
	@export ENV=$(env) && go run cmd/migrate/*.go

clean-up: .bin-validator ## Cleanup database, git branch, and cache
	@rm dealls_test.db
	@git branch | grep -v "master" | xargs git branch -D