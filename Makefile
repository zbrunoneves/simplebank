include .env

.PHONY: run
run:
	@echo "starting server..."
	@go run main.go --debug
	@echo "done."

.PHONY: test
test:
	ginkgo -r --keep-going --randomize-all --randomize-suites --fail-on-pending --trace --race -cover

.PHONY: db
db:
	@echo "starting database..."
	@docker run --name mysql --rm \
		-e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
		-e MYSQL_DATABASE=${MYSQL_DATABASE} \
		-e MYSQL_USER=${MYSQL_USER} \
		-e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
		-p 3306:3306 -d mysql:9.1.0
	@echo "done."

.PHONY: sqlc
sqlc:
	@echo "generating sql code..."
	@sqlc generate
	@echo "done."
