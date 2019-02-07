build:
	@go build

test:
	@# remove existing container
	@docker ps -q --filter "name=switch-database" | grep -q . && docker stop switch-database && docker rm -fv switch-database || true
	
	@# start new container
	@docker build -t switch-database database > /dev/null
	@docker run -d -p 5432:5432 --name switch-database switch-database > /dev/null

	@# wait for database
	@bash ../wait-for-database.sh > /dev/null

	@# run the tests
	@go test -v

	@# remove the container
	@docker kill switch-database > /dev/null
	@docker rm switch-database > /dev/null