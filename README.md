# Smart Home Adapters

## Components

### [Web Server](webserver/README.md)

Receives sensor data and sends servo commands to robots. Exposes a REST API consumed by the android app (and other integrations).

### [Android App](android/README.md)

Inteface to setup, calibrate and control the robots

### [Robot](robot/README.md)

Connects the hardware to the web server.

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