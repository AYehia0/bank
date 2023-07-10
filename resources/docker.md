# Learning about docker

## Installation
To install docker, follow the offical instructions [here](https://docs.docker.com/get-docker/)
After installing docker make sure it's up and running.

- To check running images : `docker images`
- To create a container for postgres image:

        docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

    - `--name postgres12`: Specifies the name of the container as "postgres12".
    - `-p 5432:5432`: Maps the host's port 5432 to the container's port 5432. This allows access to the PostgreSQL service running inside the container.
    - `-e POSTGRES_USER=root`: Sets the environment variable `POSTGRES_USER` to "root". This defines the username for the PostgreSQL database.
    - `-e POSTGRES_PASSWORD=secret`: Sets the environment variable `POSTGRES_PASSWORD` to "secret". This defines the password for the PostgreSQL database.
    - `-d`: Runs the container in the background (detached mode).
    - `postgres:12-alpine`: Specifies the image to use for the container, in this case, "postgres" version 12 with the "alpine" variant. The image will be pulled from Docker Hub if not already available.

- To run docker commands : `docker exec -it container_name command -U user`
In our case, `docker exec -it postgres12 psql -U root`

- Viewing logs of the container : `docker logs container_name`
