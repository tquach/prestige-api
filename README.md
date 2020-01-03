# Prestige API

All you need to run this server is Docker and Docker Machine.

* Install and set up [Docker Machine](https://docs.docker.com/machine/)
* Install [Docker](https://docker.com)

You should see this:
```
⚡ docker-machine status
Running
```

Take note of the Docker IP address by running:
```
⚡ docker-machine ip default
```

This should output the IP address. Export this as an environment variable:
```
⚡ export DOCKER_IP=$(docker-machine ip default)
```

## Running Redis locally

We are using the [Docker container for Redis 3](https://hub.docker.com/_/redis/).

```bash
⚡ docker run -d --name redis-master -p 6379:6379 redis
```

We will expose port 6379 for now but you can also link containers together.

## Running a local database

We are using the [Docker container for Postgres 9.1](https://hub.docker.com/_/postgres/).
```bash
⚡ docker run -d --rm --name postgres -p 5432:5432 -e POSTGRES_USER=prestige -e POSTGRES_PASSWORD=changeme postgres
```

We are using `-p 5432:5432` in case you want to use a database client like pgAdmin3 to connect to the database.

## Running a standalone Datadog Statsd

We use Datadog for metrics. You can run a standalone docker container for Datadog Statsd but you should use your own access token:

```bash
⚡ docker run -d --name dogstatsd -h `hostname` -v /var/run/docker.sock:/var/run/docker.sock -v /proc/mounts:/host/proc/mounts:ro \
	-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro -e DOGSTATSD_ONLY=true \
	-e API_KEY=203f4a84c721058f1b62691c75d001a3 -p 127.0.0.1:8125:8125 datadog/docker-dd-agent
```

## Database Schema

TODO

## Starting up the application server in a container

Once you have Postgres, Redis and Statsd running locally in their own containers, you can start up the application in it's own container as well:

```
⚡ make start
```

That should do everything for you, then you can access the API server through the Docker IP using curl:
```
⚡ curl $DOCKER_IP:9000/healthcheck
```

# API Documentation

TODO

# Getting started with development

For development, you will need Go.

* Install latest version of [Go](https://golang.org)

Building and running locally outside of Docker:
```
go build .
./prestige-api --hostname localhost:9000 --databaseURL postgresql://user@localhost:5433/prestige
```
