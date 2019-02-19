build:
	@go generate
	@go build -o server

test:
	@# remove existing container
	@docker ps -q --filter "name=info-database" | grep -q . && docker stop info-database && docker rm -fv info-database || true
	
	@# start new container
	@docker build -t info-database database > /dev/null
	@docker run -d -p 5432:5432 --name info-database info-database > /dev/null

	@# wait for database
	@bash ../wait-for-database.sh > /dev/null

	@# run the tests
	@go test

	@# remove the container
	@docker kill info-database > /dev/null
	@docker rm info-database > /dev/null
