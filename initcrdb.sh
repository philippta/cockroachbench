#!/bin/bash
docker compose exec roach1 ./cockroach --host=roach1:26357 init --insecure
