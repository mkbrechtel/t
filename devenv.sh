#!/bin/sh -e
    
podman build -t t-devenv --target devenv .

# Create a named volume for the home directory
if ! podman volume ls | grep -q "devenv-root-home"; then
    podman volume create devenv-root-home
fi

exec podman run -it --replace \
    --name t-devenv \
    --hostname t-devenv \
    -v claude-root-home:/root \
    -v ./:/mnt/t/ \
    --workdir /mnt/t \
    t-devenv \
    fish
