# Centralized Configuration Service
![Golang](https://img.shields.io/badge/Go-blue?logo=go&logoColor=white)
![Gorilla](https://img.shields.io/badge/Gorilla-yellow?logo=go&logoColor=white)
![Consul](https://img.shields.io/badge/Consul-pink?logo=consul&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-orange?logo=prometheus&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-blue?logo=docker&logoColor=white)
![Jaeger](https://img.shields.io/badge/Jaeger-black?logo=jaeger&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-green?logo=swagger&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-orange?logo=postman&logoColor=black)

## Technologies Used
Main technologies which are used in this project are **Golang** with **Consul** DB. For sending APIs you can use **Swagger** or **Postman** also app is containerized with **Docker**.

## Overview:  
There are two main components in this app: **Web Service** which handle requests and **Consul** DB.  
Also there are two side components: First for storing logs and traces - **Jaeger** and second for storing metrics - **Prometheus**. 

## Getting Started:  
To start this project you need to have installed next things:  
**Visual Studio Code or other IDE**, **Golang**, **Docker** and optional **Postman**.

## POSTMAN
**Postman Docs/Collection**: .............................................................treba dodati link

## SWAGGER
.............................................................................................treba ispisati nesto ukratko za swagger 2-3 recenice

## Idempotency  
**What is Idempotency middleware** ? The idempotency middleware ensures that repeated requests with the same parameters produce the same result, regardless of how many times they are sent. It helps prevent unintended side effects caused by duplicate requests, such as duplicate charges in a payment system or duplicate updates in a database. By generating and storing a unique identifier for each request and its corresponding response, the middleware can check incoming requests against this identifier. If a request with the same identifier is received again, the middleware can retrieve the previous response associated with that identifier and return it without executing the request handler again. This middleware adds an extra layer of reliability and safety to your application, especially in distributed systems where duplicate requests are more likely to occur.  
We are storing Idempotency-Key in our **Consul** DB.  

## Database:  
**Consul** is a NoSQL database designed for storing key-value pairs. We chose Consul for its simplicity and suitability for our project specifications. To access the **Consul UI**, use the port **8500**.  
This will allow you to manage and interact with your Consul data effortlessly.

## Testing:  
We have implemented unit tests for all services in this project.  
These unit tests are designed to ensure the functionality of individual components in isolation.  

## Metrics:
**Prometheus** is an open-source systems monitoring and alerting toolkit designed for reliability and scalability. It collects and stores its metrics as time-series data, providing a powerful query language called **PromQL** to query and visualize the data. You can access **Prometheus UI** on port **9090**.  

**http_total_requests** -> Total number of HTTP requests in last 24h.  
**http_successful_requests** -> Number of successful HTTP requests in last 24h (2xx, 3xx).  
**http_unsuccessful_requests** -> Number of unsuccessful HTTP requests in last 24h (4xx, 5xx).  
**average_request_duration_seconds** -> Average request duration for each endpoint.   
**requests_per_time_unit** -> Number of requests per time unit (e.g., per minute or per second) for each endpoint.  



## Tracing:  
..............................................................................treba ispisati dokumentaciju za trejsing



## Deploy:  

### Dockerfile  
**We used Multi-Stage build for lighter final image**  
**BUILD ENVIROMENT**  
**FROM golang:1.22-alpine AS build:** Use Go 1.22 on Alpine Linux as the build environment.  
**WORKDIR /app:** Set the working directory inside the container to /app.  
**COPY go.mod go.sum ./:** Copy dependency files into the container.  
**RUN go mod download:** Download Go module dependencies.  
**COPY . .:** Copy all project files into the container.  
**RUN go build -o app .:** Compile the Go application and output the binary as app.    
**RUNTIME ENVIROMENT**  
**FROM alpine:** Use the latest Alpine Linux as the runtime environment.  
**WORKDIR /app:** Set the working directory to /app.  
**COPY --from=build /app/app .:** Copy the built application binary from the build stage.  
**COPY swagger.yaml /app/swagger.yaml:** Copy the swagger.yaml file into the container.  
**EXPOSE 8000:** Expose port 8000 for the application.  
**CMD ["./app"]:** Set the command to run the application.  


### docker-compose.yml  
This docker-compose.yml file is used to define and manage the services required for app.  

**App Service**:  
-image: Specifies the Docker image to use.  
-container_name: Sets the container name.  
-hostname: Sets the hostname for the container.  
-ports: Maps a host port to a container port using an environment variable.  
-depends_on: Ensures that consul and jaeger services are started before this service.  
-networks: Connects the service to a user-defined network.  
-environment: Sets environment variables for the container.    

**Consul Service**:  
-image: Uses the Consul image.  
-ports: Maps the Consul UI port.  
-command: Runs Consul as a server with specific options.  
-volumes: Mounts a volume for persistent data storage.  
-networks: Connects the service to a user-defined network.  

**Prometheus Service**:  
-image: Uses the Prometheus image.  
-ports: Maps the Prometheus UI port.  
-volumes: Mounts directories for Prometheus configuration and data.  
-networks: Connects the service to a user-defined network.  

**Jaeger Service**:  
-image: Uses the Jaeger all-in-one image.  
-ports: Maps Jaeger ports for tracing and the Jaeger UI.  
-networks: Connects the service to a user-defined network.    

## CI PIPELINE  
.................................................................................................................... treba uraditi i ovo   

## Authors  

**Andrej Stjepanović**  
Student at the **Faculty of Technical Sciences** in Novi Sad  
Undergraduate of Software Engineering  

**Aleksander Zajac**  
Student at the **Faculty of Technical Sciences** in Novi Sad  
Undergraduate of Software Engineering   

**Dragan Bijelić**  
Student at the **Faculty of Technical Sciences** in Novi Sad  
Undergraduate of Software Engineering  
