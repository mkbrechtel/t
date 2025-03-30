#!/bin/sh -e
    
podman build -t t5-devenv --target devenv .

# Create a named volume for the home directory
if ! podman volume ls | grep -q "devenv-root-home"; then
    podman volume create devenv-root-home
fi

exec podman run -it --replace \
    --name t5-devenv \
    --hostname t5-devenv \
    -v claude-root-home:/root \
    -v ./:/mnt/t5/ \
    --workdir /mnt/t5 \
    t5-devenv \
    fish
