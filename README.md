# Golang Microservice

## Introduction

### Hey GO-pher

This is a basic microservice architecture written in Golang to strengthen my grip on microservices and Golang.

## Featuring

Currently, it features authentication services, with plans to add more services in the future.

## Setup Locally

### Using Go

To set this up locally, follow these steps:

1. Download all of the required dependencies:

```bash
go mod download
```

2. Build the project:

```bash
go build main.go
```

3. Run the project:

```bash
go run main.go
```

### Using Docker

You can also run the project using Docker. Follow the steps below:

1. Pull the pre-built Docker image from Docker Hub:

```bash
docker pull deshwalankush23/microservice
```

2. Run the Docker container:

```bash
docker run -d -p 3000:3000 deshwalankush23/microservice
```

This command will run the container and map port 3000 of the container to port 3000 on your local machine.

### Using Docker Compose

Alternatively, if you have Docker Compose installed, you can use the following command to start the application:

```bash
docker-compose up
```

This will start the service based on the docker-compose.yml configuration file.

## Enjoy Coding with Golang!