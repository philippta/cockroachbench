#!/bin/bash

apt update && apt install -y postgresql vim

# Create user for replication
sudo -u postgres psql -c 'create user repuser replication'

# Allow connections on all interfaces
echo "listen_addresses = '*'" >> /etc/postgresql/15/main/postgresql.conf

# Set replication mode to synchronous
echo "synchronous_commit = on" >> /etc/postgresql/15/main/postgresql.conf

# Ensure all replicas are synchronous replicas
echo "synchronous_standby_names = '*'" >> /etc/postgresql/15/main/postgresql.conf

# Allow benchmark client and replication user to connect
echo "host all postgres 10.0.0.0/24 trust" >> /etc/postgresql/15/main/pg_hba.conf
echo "host all repuser 10.0.0.0/24 trust" >> /etc/postgresql/15/main/pg_hba.conf
echo "host replication repuser 10.0.0.0/24 trust" >> /etc/postgresql/15/main/pg_hba.conf

# Restart postgres to reflect the changes
systemctl restart postgresql

# Print some instructions to run the benchmark
ip=$(hostname -I | cut -d" " -f2)
echo "====================================================================="
echo "Set up the standby database now. And after that..."
echo "---------------------------------------------------------------------"
echo "Run run the benchmark execute these commands on client:"
echo "---------------------------------------------------------------------"
echo "pgbench -i -d postgres -h $ip -p 5432 -U postgres -n -s 100"
echo "pgbench postgres -h $ip -p 5432 -U postgres -n -t 5000"
echo "====================================================================="

# How to connect to the database locally
# sudo -u postgres psql
