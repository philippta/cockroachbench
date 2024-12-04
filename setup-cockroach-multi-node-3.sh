#!/bin/bash

IP_NODE_1=10.0.0.3
IP_NODE_2=10.0.0.5
IP_NODE_3=10.0.0.2

# Download and extract cockroach binaries
wget https://binaries.cockroachdb.com/cockroach-v24.3.0.linux-amd64.tgz
tar xf cockroach-v24.3.0.linux-amd64.tgz

# Copy binaries to better directory
cp -i cockroach-v24.3.0.linux-amd64/cockroach /usr/local/bin/

# No idea what this geos thing is, but it says to in the guide
mkdir -p /usr/local/lib/cockroach
cp -i cockroach-v24.3.0.linux-amd64/lib/libgeos.so /usr/local/lib/cockroach/
cp -i cockroach-v24.3.0.linux-amd64/lib/libgeos_c.so /usr/local/lib/cockroach/

# Create cockroach user
useradd cockroach

# Create data directory
mkdir /var/lib/cockroach
chown -R cockroach /var/lib/cockroach

# Get systemd service files and update them for single node cluster
# Note: --cache=.25 --max-sql-memory=.25 is a recommendation from Cockroach Labs
curl -o /etc/systemd/system/insecurecockroachdb.service https://raw.githubusercontent.com/cockroachdb/docs/main/src/current/_includes/v24.3/prod-deployment/insecurecockroachdb.service
sed -i "s|ExecStart=/usr/local/bin/cockroach start --insecure --advertise-addr=<node1 address> --join=<node1 address>,<node2 address>,<node3 address> --cache=.25 --max-sql-memory=.25|ExecStart=/usr/local/bin/cockroach start --insecure --advertise-addr=$IP_NODE_3 --join=$IP_NODE_1,$IP_NODE_2,$IP_NODE_3 --cache=.25 --max-sql-memory=.25|" /etc/systemd/system/insecurecockroachdb.service

# Start cockroachdb service
systemctl daemon-reload
systemctl start insecurecockroachdb

# Init the cluster
cockroach init --insecure --host=10.0.0.2

# Print some instructions to run the benchmark
ip=$(hostname -I | cut -d" " -f2)
echo "====================================================================="
echo "Run run the benchmark execute these commands on client:"
echo "---------------------------------------------------------------------"
echo "pgbench -i -d defaultdb -h $ip -p 26257 -U root -n -s 100"
echo "pgbench defaultdb -h $ip -p 26257 -U root -n -t 5000"
echo "====================================================================="
