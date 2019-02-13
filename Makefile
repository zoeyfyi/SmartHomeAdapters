#
# CHECK
#

check-docker-deps:
	@which docker > /dev/null
	@which docker-compose > /dev/null

check-go-deps:
	@which protoc > /dev/null
	@which protoc-gen-go > /dev/null
	@which gometalinter > /dev/null

check-arduino-deps:
	@which arduino > /dev/null

check: check-docker-deps check-go-deps check-arduino-deps

#
# BUILD
#

build-clientserver: check-go-deps
	@(cd clientserver && go generate)
	@(cd clientserver && go build -o ../build/clientserver)

build-infoserver: check-go-deps
	@(cd infoserver && go generate)
	@(cd infoserver && go build -o ../build/infoserver)

build-robotserver: check-go-deps
	@(cd robotserver && go generate)
	@(cd robotserver && go build -o ../build/robotserver)

build-switchserver: check-go-deps
	@(cd switchserver && go generate)
	@(cd switchserver && go build -o ../build/switchserver)

build-userserver: check-go-deps
	@(cd userserver && go generate)
	@(cd userserver && go build -o ../build/userserver)

build-android:
	@(cd android && ./gradlew assembleDebug)
	@cp android/app/build/outputs/apk/debug/app-debug.apk build/app-debug.apk

build: build-clientserver build-infoserver build-robotserver build-switchserver build-userserver build-android

#
# DOCKER
#

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

#
# CLEAN
#

clean:
	@rm -rf build/*

#
# LINT
#

lint-clientserver:
	@(cd clientserver && gometalinter --disable=gocyclo --enable=gofmt ./...)

lint-infoserver:
	@(cd infoserver && gometalinter --disable=gocyclo --enable=gofmt ./...)

lint-robotserver:
	@(cd robotserver && gometalinter --disable=gocyclo --enable=gofmt ./...)

lint-switchserver:
	@(cd switchserver && gometalinter --disable=gocyclo --enable=gofmt ./...)

lint-userserver:
	@(cd userserver && gometalinter --disable=gocyclo --enable=gofmt ./...)

lint-android:
	@(cd android && ./gradlew lint)

lint: lint-clientserver lint-infoserver lint-robotserver lint-switchserver lint-userserver lint-android
