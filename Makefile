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
	@docker build -f go.Dockerfile -t smarthomeadapters/clientserver . --build-arg SERVICE=clientserver

docker-infodb:
	@(cd infodb && docker build -t smarthomeadapters/infodb .)

docker-infoserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/infoserver . --build-arg SERVICE=infoserver

docker-robotserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/robotserver . --build-arg SERVICE=robotserver

docker-switchdb:
	@(cd switchdb && docker build -t smarthomeadapters/switchdb .)

docker-switchserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/switchserver . --build-arg SERVICE=switchserver

docker-userdb:
	@(cd userdb && docker build -t smarthomeadapters/userdb .)

docker-userserver:
	@docker build -f go.Dockerfile -t smarthomeadapters/switchserver . --build-arg SERVICE=userserver

docker: docker-clientserver docker-infoserver docker-robotserver docker-switchserver docker-userserver docker-infodb docker-switchdb docker-userdb

docker-push:
	@docker tag smarthomeadapters/clientserver smarthomeadapters/clientserver:latest
	@docker push smarthomeadapters/clientserver:latest
	@docker tag smarthomeadapters/infodb smarthomeadapters/infodb:latest
	@docker push smarthomeadapters/infodb:latest
	@docker tag smarthomeadapters/infoserver smarthomeadapters/infoserver:latest
	@docker push smarthomeadapters/infoserver:latest
	@docker tag smarthomeadapters/robotserver smarthomeadapters/robotserver:latest
	@docker push smarthomeadapters/robotserver:latest
	@docker tag smarthomeadapters/switchdb smarthomeadapters/switchdb:latest
	@docker push smarthomeadapters/switchdb:latest
	@docker tag smarthomeadapters/switchserver smarthomeadapters/switchserver:latest
	@docker push smarthomeadapters/switchserver:latest
	@docker tag smarthomeadapters/userdb smarthomeadapters/userdb:latest
	@docker push smarthomeadapters/userdb:latest

docker-push-test:
	@docker tag smarthomeadapters/clientserver smarthomeadapters/clientserver:test
	@docker push smarthomeadapters/clientserver:test
	@docker tag smarthomeadapters/infodb smarthomeadapters/infodb:test
	@docker push smarthomeadapters/infodb:test
	@docker tag smarthomeadapters/infoserver smarthomeadapters/infoserver:test
	@docker push smarthomeadapters/infoserver:test
	@docker tag smarthomeadapters/robotserver smarthomeadapters/robotserver:test
	@docker push smarthomeadapters/robotserver:test
	@docker tag smarthomeadapters/switchdb smarthomeadapters/switchdb:test
	@docker push smarthomeadapters/switchdb:test
	@docker tag smarthomeadapters/switchserver smarthomeadapters/switchserver:test
	@docker push smarthomeadapters/switchserver:test
	@docker tag smarthomeadapters/userdb smarthomeadapters/userdb:test
	@docker push smarthomeadapters/userdb:test

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

lint-docker-compose:
	docker-compose config

lint: lint-clientserver lint-infoserver lint-robotserver lint-switchserver lint-userserver lint-android lint-docker-compose

#
# TEST
#

test-clientserver:
	@(cd clientserver && go test)

test-infoserver:
	@(cd infoserver && go test)

test-robotserver:
	@(cd robotserver && go test)

test-switchserver:
	@(cd switchserver && go test)

test-userserver:
	@(cd userserver && go test)

test-android:
	@(cd android && ./gradlew test)

test-services: test-clientserver test-infoserver test-robotserver test-switchserver test-userserver

test: test-services test-android

#
# Reports
#

compile-reports:
	for DIR in reports/*/; \
	do \
		echo "Compiling $${DIR}"; \
		docker run --mount src=$$PWD/$${DIR},target=/usr/src/tex,type=bind dxjoke/tectonic-docker /bin/sh -c "tectonic document.tex"; \
	done

#
# CI
#

ci-test-android: 
	@docker run -it --rm -v $$PWD/android:/root/tmp budtmo/docker-android-x86-9.0 bash -c "(cd tmp && ./gradlew test)"

ci: docker test-services ci-test-android compile-reports lint-docker-compose