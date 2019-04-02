# Smart Home Adapters

## Git

### Commit messages

Follow the [conventional commits spec](https://www.conventionalcommits.org/en/v1.0.0-beta.2/).

### Making changes

#### Hot fixes & documentation

This is strictly for hot fixes and non-code changes.

1. Branch from master 
2. Create a pull request
3. Wait for merge

#### Features

For all other code changes

1. Branch from the relevent component branch
2. Create a pull request (into component branch)
3. Have someone else review your code 
4. Wait for merge
5. ~3 days before a demo all the component branches will be merged into master

## Development

- `make check` check build dependencies are installed
- `make build` to build everything
    - `make build-android`
    - `make build-clientserver`
    - `make build-infoserver`
    - ...
- `make docker` builds all the docker images
    - `make docker-push` push docker images (tagged latest)
    - `make docker-push-test` push docker images (tagged test)
- `make clean` deletes the build folder
- `make lint` runs various linters
    - `make lint-docker-compose` checks docker compose file is valid
    - `make lint-android`
    - `make lint-clientserver`
    - `make lint-robotserver`
    - ...
- `make test` test all the projects
    - `make test-android`
    - `make test-clientserver`
    - `make test-infoserver`
    - ...
- `make compile-reports` builds all the reports

There are other commands, check the `Makefile`.
