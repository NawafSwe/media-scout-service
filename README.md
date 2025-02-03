# Media Scout Service

Backend service fetches media information through ITunes API.

## Tech Stack

The Media Scout Service project uses the following tech stack:

- **Programming Language:** Go
- **Frameworks and Libraries:**
    - `net/http` for HTTP server
    - `sqlx` for database interactions
    - `kit` for endpoint and transport layers
    - `gomock` for generating mocks
- **Database:** PostgreSQL
- **Containerization:** Docker and Docker Compose
- **Build and Dependency Management:** Go Modules (`go.mod`)
- **Testing:** Go's testing package with unit and integration tests
- **Logging:** Custom logging interface with context-based logging
- **Environment Configuration:** `.env` file for environment variables
- **Makefile:** For build, run, and management tasks


## Prerequisites

Before you begin, ensure you have the following installed on your local machine:

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [PostgreSQL](https://www.postgresql.org/download/)

## Getting Started

### Clone the Repository

Clone the repository to your local machine using the following command:

```sh
git clone https://github.com/NawafSwe/media-scout-service.git
cd media-scout-service
```

## Environment Variables
Create a .env file in the root directory of the project and add the following environment variables:
```.dotenv
# GENERAL CONFIG
TLS_ENABLED=false
LOGGING_ENABLED=true
LOGGING_LEVEL=debug

# HTTP CONFIG
HTTP__PORT=3001

# DB CONFIG
DB__DSN=postgres://postgres:1234@postgres:5432/media_scout_db?sslmode=disable
```

## Running the Project
Using Docker
Build and Run Docker Containers:  
``` docker-compose up --build -d ```
This command will build the Docker images and start the containers for the application and PostgreSQL database.
Then run ```make build-app && make http``` to spin the http server.
Access the Application:  The application will be running at http://localhost:3001.


Without Docker
Start PostgreSQL:  Ensure PostgreSQL is running and accessible with the credentials provided in the .env file.  
Run Database Migrations:  Apply the necessary database migrations (if any).  
Install Dependencies:  
```go mod tidy```
Run the Application:  
```go run cmd/main.go```
The application will be running at http://localhost:3001. 


## API Endpoints

### Health Check

- **URL:** `/health`
- **Method:** `GET`
- **Description:** Checks the health of the service.

### Search Media

- **URL:** `/api/v1/media/search`
- **Method:** `GET`
- **Query Parameters:**
    - `term` (string): The search term.
    - `limit` (int, optional): The number of results to return (default is 20).
- **Description:** Searches for media information using the iTunes API.