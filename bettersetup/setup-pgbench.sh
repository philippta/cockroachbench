#!/bin/bash

apt update && apt install -y postgresql vim

pgbench -i -F 1000 -d postgres -h 10.0.0.2 -p 5432 -n
