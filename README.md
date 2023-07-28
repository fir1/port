## About
This repository demonstrates the use of Domain-Driven Design (DDD) and hexagonal architecture to build a scalable and efficient REST API. 
The main feature of this API is its ability to handle very large files with an unknown size, as reading such files at once could exceed 
the limited RAM of computers (e.g., 200MB). To address this, the application implements a streaming-based approach to read and process 
large JSON files in smaller chunks, reducing memory usage.

## How to run program locally 
To run the program locally, use the following command:

````
go run cmd/*.go 
````

## How to run linters
To ensure code quality and consistency, it is recommended to use a linter for the project. 
Before running the linter, make sure you have it installed on your computer. 
You can follow the installation instructions provided at https://golangci-lint.run/usage/install/ to set up the linter locally.

Once the linter is installed, you can run the following command to analyze the entire project for potential issues:
``golangci-lint run``

Executing this command will trigger the linter to examine the codebase, identifying any code style violations, potential bugs, 
or other issues that may need attention. It is essential to address any reported problems to maintain code quality and adhere 
to best practices. Regularly running the linter during development can help catch and fix problems early in the development 
process, leading to a more maintainable and robust codebase.

## Run tests
To run the tests, use the command:
`go test ./...`

## API Endpoints
The application currently provides the following endpoints:
1. ``GET /health``: This endpoint checks the health of the server.

2. ``POST /port``: This endpoint creates ports and saves them into the database using the default `ports.json` file. 
The streaming-based approach is used to efficiently handle large JSON files.

3. ``POST /ports/from-file``: The service will parse the file and save the ports into the database. This API also handles very large files as chunks, ensuring efficient processing.
The request `Content-Type: multipart/form-data` to send files to server.

## Features
- The API server utilizes caching to improve response times. Currently, memory caching from the github.com/allegro/bigcache/v3 library is used, but it can be replaced with other caching solutions, such as redis, by implementing the `CacheClientInterface` in `pkg/cache/cache.go`. The use of interfaces allows for easy swapping of caching implementations without changing the application details.

- The project benefits from automatic dependency injection tools, such as `"uber/fx"`, to manage dependencies and facilitate modular and testable code.

- The project uses linters, such as `golangci`, to ensure code quality

- Swagger documentation is implemented for the API, providing better API visibility and documentation

- Test coverage is provided wherever feasible to ensure code quality.

- The application is dockerized, making it easy to deploy and run in containers.

- Use of clean code with domain driven design and hexagonal architecture

## SWAGGER Documentation
The application also has SWAGGER documentation that provides detailed information about the API endpoints. To access the documentation, run the server using the command `go run cmd/*.go` and visit http://localhost:8080/swagger/index.html in your browser.


## Docker

In this section, we will explain how to create a Docker image, run the web service, and publish the Docker image to Docker Hub.

**Creating an image from the Dockerfile**
````
docker build -t 2112fir/port -f build/Dockerfile .
````

**Publishing to Docker Hub**
To publish the image to Docker Hub, you need to log in first using either of the following methods:

Login to Docker Hub via CLI
1) Direct usage of password in CLI is not recommended 
````
docker login -p {Passowrd} -u {Username}
````

2) Creating an access token from https://hub.docker.com/settings/security (preferred way):
````
docker login -u 2112fir
````

At the password prompt, enter the personal access token.

You can push a new image to this repository using the CLI
````
docker push 2112fir/port
````

**Running server from the above created docker image**
````
docker run -p 8080:8080 2112fir/port
````

Later on we will use publicly pushed image inside Kubernetes manifest.

Alternatively, you can use the Docker Compose file to build the image locally and test it from your local Docker environment:
````
cd build

docker-compose up
````


# What can be improved ?
- Currently, the BDD (Behavior-Driven Development) integration tests are missing from the project due to time constraints for this project. In real projects, the use of BDD tests with the help of the Cucumber framework would be useful to check the behavior of the business logic, as it allows making client REST API calls and verifying the expected behavior.