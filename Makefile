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

build-android:
	@(cd android && ./gradlew assembleDebug)

#
# DOCKER
#


docker: docker-compose build

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
	@docker tag smarthomeadapters/userserver smarthomeadapters/userserver:latest
	@docker push smarthomeadapters/userserver:latest
	@docker tag smarthomeadapters/thermodb smarthomeadapters/thermodb:latest
	@docker push smarthomeadapters/thermodb:latest
	@docker tag smarthomeadapters/thermostatserver smarthomeadapters/thermostatserver:latest
	@docker push smarthomeadapters/thermostatserver:latest

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
	@docker tag smarthomeadapters/userserver smarthomeadapters/userserver:test
	@docker push smarthomeadapters/userserver:test
	@docker tag smarthomeadapters/thermodb smarthomeadapters/thermodb:test
	@docker push smarthomeadapters/thermodb:test
	@docker tag smarthomeadapters/thermostatserver smarthomeadapters/thermostatserver:test
	@docker push smarthomeadapters/thermostatserver:test
	@docker tag smarthomeadapters/account-app smarthomeadapters/account-app:test
	@docker push smarthomeadapters/account-app:test

#
# LINT
#

lint-services:
	@for d in services/*/ ; do (cd "$d" && test -f go.mod && golangci-lint run); done

lint-android:
	@(cd android && ./gradlew lint)

lint-docker-compose:
	docker-compose config

lint: lint-services lint-android lint-docker-compose

#
# DB
#

db-up:
	@docker run --rm -d --name test_infodb -p 5001:5432 -e POSTGRES_USER=temp -e POSTGRES_PASSWORD=password smarthomeadapters/infodb
	@docker run --rm -d --name test_switchdb -p 5002:5432 -e POSTGRES_USER=temp -e POSTGRES_PASSWORD=password smarthomeadapters/switchdb
	@docker run --rm -d --name test_thermodb -p 5003:5432 -e POSTGRES_USER=temp -e POSTGRES_PASSWORD=password smarthomeadapters/thermodb
	@docker run --rm -d --name test_userdb -p 5004:5432 -e POSTGRES_USER=temp -e POSTGRES_PASSWORD=password smarthomeadapters/userdb

db-down:
	@docker stop test_infodb
	@docker stop test_switchdb
	@docker stop test_thermodb
	@docker stop test_userdb

#
# TEST
#

test-android:
	@(cd android && ./gradlew test)

test-services:
	@for d in services/*/ ; do (cd "$d" && test -f go.mod && go test); done

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

ci: docker-dbs test-services ci-test-android compile-reports lint-docker-compose
