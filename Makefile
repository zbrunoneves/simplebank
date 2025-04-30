include .env

.PHONY: run
run:
	@go run main.go --debug

.PHONY: test
test:
	ginkgo -r --keep-going --randomize-all --randomize-suites --fail-on-pending --trace --race -cover

.PHONY: db
db:
	@echo "starting database..."
	@docker run --name mysql --rm \
		-e MYSQL_ROOT_PASSWORD=123456 \
		-e MYSQL_DATABASE=bank \
		-e MYSQL_USER=admin \
		-e MYSQL_PASSWORD=123 \
		-p 3306:3306 -d mysql:9.3.0
	@echo "done."