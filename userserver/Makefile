build:
	@go build

test:
	@# remove existing container
	@docker ps -q --filter "name=user-database" | grep -q . && docker stop user-database && docker rm -fv user-database || true
	
	@# start new container
	@docker build -t user-database database > /dev/null
	@docker run -d -p 5432:5432 --name user-database user-database > /dev/null

	@# wait for database
	@bash ../wait-for-database.sh > /dev/null

	@# run the tests
	@go test

	@# remove the container
	@docker kill user-database > /dev/null
	@docker rm user-database > /dev/null