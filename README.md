# PostgreSQL vs. CockroachDB Benchmark

This benchmarking report compares the performance of **PostgreSQL** and **CockroachDB** using the **TPC-B** benchmark. 
TPC-B simulates a banking environment with high transaction throughput, focusing on operations like account balance updates, teller transactions, and branch fund management.

What I am benchmarking:

- PostgreSQL Single Instance
- PostgreSQL Primary/Standby with Synchronous Streaming Replication
- CockroachDB Single Node Cluster
- CockroachDB 3 Node Cluster

## Environment

All databases were running on **CCX13** instances from Hetzner, with these specs:

- 2 vCPUs (AMD, dedicated)
- 8 GB RAM
- 80 GB SSD

The benchmark tool was running on a **CX22** instance from Hetzner, with these specs:

- 2 vCPUs (Intel/AMD, shared)
- 4 GB RAM
- 40 GB SSD

For the CockroachDB multi node setup or the PostgreSQL primary/standby setup, each node/instance was running on it's own dedicated server on the same local network.

The benchmark was run from another seperate server, also within the same local network, to not have any impact on database performance.

### OS and Versions

- Debian 12
- CockroachDB 23.4
- PostgreSQL 15.10

## Benchmark Setup

In this setup I'm using the `pgbench` utility provided by a default PostgreSQL installation. 
It runs the TPC-B benchmark as mentioned above.

### Seeding the database

```
pgbench -i -d <database> -h <host> -p <port> -U <user> -n -s 100
```

It was configured with a scale factor of **100**.
This means the following amounts of records are created:

- 10,000,000 Accounts
- 1,000 Tellers
- 100 Branches

### Running the benchmark

The following command simulates **5000** transactions against the database.

```
pgbench <database> -h <host> -p <port> -U <user> -n -t 5000
```


## Benchmark Results

### PostgreSQL Single Instance 

```
pgbench (15.10 (Debian 15.10-0+deb12u1))
transaction type: <builtin: TPC-B (sort of)>
scaling factor: 100
query mode: simple
number of clients: 1
number of threads: 1
maximum number of tries: 1
number of transactions per client: 5000
number of transactions actually processed: 5000/5000
number of failed transactions: 0 (0.000%)
latency average = 2.530 ms
initial connection time = 7.909 ms
tps = 395.250137 (without initial connection time)
```

### PostgreSQL Primary/Standby with Synchronous Streaming Replication

```
pgbench (15.10 (Debian 15.10-0+deb12u1))
transaction type: <builtin: TPC-B (sort of)>
scaling factor: 100
query mode: simple
number of clients: 1
number of threads: 1
maximum number of tries: 1
number of transactions per client: 5000
number of transactions actually processed: 5000/5000
number of failed transactions: 0 (0.000%)
latency average = 3.647 ms
initial connection time = 8.228 ms
tps = 274.221927 (without initial connection time)
```

### CockroachDB Single Node Cluster

```
pgbench (15.10 (Debian 15.10-0+deb12u1), server 13.0.0)
transaction type: <builtin: TPC-B (sort of)>
scaling factor: 100
query mode: simple
number of clients: 1
number of threads: 1
maximum number of tries: 1
number of transactions per client: 5000
number of transactions actually processed: 5000/5000
number of failed transactions: 0 (0.000%)
latency average = 6.280 ms
initial connection time = 1.440 ms
tps = 159.244610 (without initial connection time)
```

### CockroachDB 3 Node Cluster

```
pgbench (15.10 (Debian 15.10-0+deb12u1), server 13.0.0)
transaction type: <builtin: TPC-B (sort of)>
scaling factor: 100
query mode: simple
number of clients: 1
number of threads: 1
maximum number of tries: 1
number of transactions per client: 5000
number of transactions actually processed: 5000/5000
number of failed transactions: 0 (0.000%)
latency average = 10.490 ms
initial connection time = 1.668 ms
tps = 95.329922 (without initial connection time)
```

## Benchmark Summary

| Database Configuration          | Transactions per Second | Difference              |
| ------------------------------- | ----------------------- | ----------------------- |
| PostgreSQL Single Instance      | 395.250137              | 100.00%                 |
| PostgreSQL Primary/Standby      | 274.221927              | 69.37% (1.44x slowdown) |
| CockroachDB Single Node Cluster | 159.244610              | 40.28% (2.48x slowdown) |
| CockroachDB 3 Node Cluster      | 95.329922               | 24.11% (4.15x slowdown) |
