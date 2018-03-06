#!/usr/bin/env bash
# Check versions
docker version
docker-compose version
git version

# Checkout the branch
git clone https://github.com/jsotogaviard/recipes
cd recipes
git checkout master

# Start docker
docker-compose up -d --build