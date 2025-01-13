include .env

.PHONY: db
db:
	@echo "starting database..."
	@docker run --name mysql --rm \
		-e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
		-e MYSQL_DATABASE=${MYSQL_DATABASE} \
		-e MYSQL_USER=${MYSQL_USER} \
		-e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
		-p 3306:3306 -d mysql:9.1.0
	@echo "database started."

.PHONY: sqlc
sqlc:
	@echo "generating sql code..."
	@sqlc generate
	@echo "sql code generated."
