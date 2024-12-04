#!/bin/bash

apt update && apt install -y postgresql

# pgbench <database> -h <host> -p <port> -U <user> -n -t 5000
