# FampayAssignment

## Update enviornment variables

    update enviornment variables in .env from .env.example file
    update docker-compose.yml files enviornment variables.

## Run Docker compose

    docker-compose up --build

## Test the API's

    # 1. All Videos (GET)
        URL : http://localhost:8080/api/v1/videos?pagenum=1&pagesize=1
    # 2. Search Video (GET)
        URL : http://localhost:8080/api/v1/search?query=live+t20+Nepal
