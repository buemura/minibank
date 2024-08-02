#!/usr/bin/bash

cd docker;
docker-compose down && docker-compose build && docker-compose up -d