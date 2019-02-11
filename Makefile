SERVERS = clientserver infoserver robotserver switchserver userserver

# Check docker dependencys
check-docker-deps:
	@which docker > /dev/null
	@which docker-compose > /dev/null

# Check go dependencys
check-go-deps:
	@which protoc > /dev/null
	@which protoc-gen-go > /dev/null

# Check arduino dependencys
check-arduino-deps:
	@which arduino > /dev/null

# Check all dependencys
check: check-docker-deps check-go-deps check-arduino-deps

# Vendors all go dependencies
vendor:
	@ for SERVER in $(SERVERS); do (cd $$SERVER && go mod vendor); done

# Builds all the servers
build-go: check-go-deps
	@ for SERVER in $(SERVERS); do (cd $$SERVER && go generate && go build -o ../build/$$SERVER); done

# Builds everything
build: build-go

# Cleans the build folder
clean:
	@rm -rf build/*