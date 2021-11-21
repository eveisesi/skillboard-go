#!/bin/bash

# I build a custom redis image for us.
# The image executes a custom init script that fetches a password from the env
# that is injected by docker compose, sed the redis.conf file and replace a placeholder
# password with the real password. This allows us to store the redis.conf in version control
# without leaking the password on Github

docker build -f ./.config/redis/scripts/Dockerfile ./.config/redis/scripts -t skillz-redis:latest

echo $GITHUB_PAT | docker login ghcr.io -u ddouglas --password-stdin