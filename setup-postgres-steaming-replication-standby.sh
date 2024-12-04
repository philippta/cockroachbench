#!/bin/bash

PRIMARY_IP=10.0.0.3

apt update && apt install -y postgresql vim

# Allow connections on all interfaces
echo "listen_addresses = '*'" >> /etc/postgresql/15/main/postgresql.conf

# Allow benchmark client and replication user to connect
echo "host all postgres 10.0.0.0/24 trust" >> /etc/postgresql/15/main/pg_hba.conf

# Stop for basebackup
systemctl stop postgresql

# Remove default database files
rm -rf /var/lib/postgresql/15/main

# Sync over database files from primary
sudo -u postgres pg_basebackup -h $PRIMARY_IP -U repuser --checkpoint=fast -D /var/lib/postgresql/15/main -R --slot=replica -C

# Start postgres again
systemctl start postgresql

# How to connect to the database locally
# sudo -u postgres psql
