# Docker Swarm Deployment Guide

This file documents the Docker Swarm orchestration requirement.

## Build Images

Docker Swarm does not build images from `build:` sections during `docker stack deploy`, so build the local images first:

```bash
docker compose build
```

## Start Swarm

```bash
docker swarm init
```

If the node is already in a swarm, continue with the deploy command.

## Deploy Stack

Make sure `.env` exists and contains the required variables, then run:

```bash
docker stack deploy -c docker-stack.yml clothes
```

## Verify

```bash
docker stack services clothes
docker stack ps clothes
docker service logs clothes_notification-service --tail 50
```

Evidence screenshots for the PDF:

- `docker node ls`
- `docker stack services clothes`
- `docker stack ps clothes`
- `docker service logs clothes_notification-service --tail 50`

## Remove Stack

```bash
docker stack rm clothes
```
