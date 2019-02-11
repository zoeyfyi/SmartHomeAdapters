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

# Builds all the servers
build-go: check-go-deps
	@ for SERVER in $(SERVERS); do (cd $$SERVER && go generate && go build -o ../build/$$SERVER); done

# Builds the android app
build-android:
	@(cd android && ./gradlew assembleDebug)
	@cp android/app/build/outputs/apk/debug/app-debug.apk build/app-debug.apk

# Builds everything
build: build-go build-android

docker-clientserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/clientserver clientserver

docker-infoserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/infoserver infoserver

docker-robotserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/robotserver robotserver

docker-switchserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/switchserver switchserver

docker-userserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/userserver userserver

docker: docker-clientserver docker-infoserver docker-robotserver docker-switchserver docker-userserver

# Cleans the build folder
clean:
	@rm -rf build/*