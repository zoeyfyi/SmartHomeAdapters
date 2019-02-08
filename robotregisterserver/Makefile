build:
	@go build

test:
	@# remove existing container
	@docker ps -q --filter "name=robot-register-database" | grep -q . && docker stop robot-register-database && docker rm -fv robot-register-database || true
	
	@# start new container
	@docker build -t robot-register-database database > /dev/null
	@docker run -d -p 5432:5432 --name robot-register-database robot-register-database > /dev/null

	@# wait for database
	@bash ../wait-for-database.sh > /dev/null

	@# run the tests
	@go test

	@# remove the container
	@docker kill robot-register-database > /dev/null
	@docker rm robot-register-database > /dev/null
