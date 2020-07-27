# point-system

A simple user's point management system.

## Features

- Save and retrieve user point
- Add user point
- Retrieve user point

## Flow

![flow](https://github.com/rl404/point-system/raw/master/flow.png "Flow")

1. User sends request to add/subtract point to REST API.
2. After validating the request, REST API will send the request to RabbitMQ queue to be processed later. User will get `202 Accepted` response.
3. Worker in the background will consume RabbitMQQ queue and process the request (add/subtract point) and update to database.


## Requirements

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Docker compose](https://docs.docker.com/compose/)

## Quick Start

1. `git clone https://github.com/rl404/point-system.git`
2. `cd point-system`
3. `make`
4. Wait until `server listen at :31001`
5. Endpoints are ready to use.

#### See worker logs

```
make docker-worker-logs
```

#### Remove containers

```
make docker-stop
make docker-rm
```

## Endpoints

- `GET` - `/v1/point`
- `POST` - `/v1/point/add`
- `POST` - `/v1/point/subtract`

*import the `postman_collection.json` for more details.*

## Tech Stacks

### [Go](https://golang.org/)

Using golang as main programming language. Go is good for creating a simple service such as this since the compiled binary is very small (around 10 MB) which can be used to create a small docker container. Won't be using framework since this is a simple service and pretty sure won't be using all of the framework features. Using framework may also affect to compiling time and compiled binary size.

### [PostgreSQL](https://postgresql.org/)

Database to keep the user point data. Since we already know the column needed to keep the data, we better use relational database.

Tables:
- `user_point` to keep user's current point.
- `log` to keep logging of user's point changes.

### [RabbitMQ](https://www.rabbitmq.com/)

Message broker to help queueing transaction. Handling transaction when the traffic/request is higher than the database speed can handle. Let the transaction processed asynchronously in the back without blocking HTTP request.

### [Docker](https://www.docker.com/) & [Docker compose](https://docs.docker.com/compose/)

Quick and easy container management and deployment without installing the stack directly to local/server.