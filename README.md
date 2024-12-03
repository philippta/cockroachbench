# PostgreSQL vs. CockroachDB (Single Node and 3-Node Cluster)

This benchmarking report compares the performance of **PostgreSQL** and **CockroachDB** (in both **single-node** and **3-node cluster** configurations) using the **TPC-B** benchmark. TPC-B simulates a banking environment with high transaction throughput, focusing on operations like account balance updates, teller transactions, and branch fund management.

## Environment

- Apple MacBook with Mac M1 Pro
- PostgreSQL and CockroachDB running inside Docker
- Docker Resources: 8 CPUs, 7.90 GB Memory

## Benchmark Results

```
====================================
POSTGRES
------------------------------------
Seeded tables with:
- 5000 accounts
- 500 branches
- 500 tellers
Benchmarking TPC-B...
Transactions:     5000
Duration:         31.723867416s
Transactions/Sec: 157.61003960942793
====================================

====================================
COCKROACH SINGLE NODE
------------------------------------
Seeded tables with:
- 5000 accounts
- 500 branches
- 500 tellers
Benchmarking TPC-B...
Transactions:     5000
Duration:         1m23.00112225s
Transactions/Sec: 60.24014934328193
====================================

====================================
COCKROACH 3 NODES
------------------------------------
Seeded tables with:
- 5000 accounts
- 500 branches
- 500 tellers
Benchmarking TPC-B...
Transactions:     5000
Duration:         2m43.683688792s
Transactions/Sec: 30.546721160186696
====================================
```

## Schema

```sql
CREATE TABLE pgbench_accounts (
	aid        INTEGER NOT NULL,
	bid        INTEGER NOT NULL,
	abalance   INTEGER NOT NULL,
	filler     CHAR(84) NOT NULL
);
CREATE TABLE pgbench_branches (
	bid        INTEGER NOT NULL,
	bbalance   INTEGER NOT NULL,
	filler     CHAR(88) NOT NULL
);
CREATE TABLE pgbench_tellers (
	tid        INTEGER NOT NULL,
	bid        INTEGER NOT NULL,
	tbalance   INTEGER NOT NULL,
	filler     CHAR(84) NOT NULL
);
CREATE TABLE pgbench_history (
	tid        INTEGER NOT NULL,
	bid        INTEGER NOT NULL,
	aid        INTEGER NOT NULL,
	delta      INTEGER NOT NULL,
	mtime      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Transaction

```sql
BEGIN;

-- Update the account balance
UPDATE pgbench_accounts
SET abalance = abalance + $1
WHERE aid = $2;

-- Update the teller balance
UPDATE pgbench_tellers
SET tbalance = tbalance + $1
WHERE tid = $2;

-- Update the branch balance
UPDATE pgbench_branches
SET bbalance = bbalance + $1
WHERE bid = $2;

-- Insert a record into the history table
INSERT INTO pgbench_history (tid, bid, aid, delta, mtime)
VALUES ($1, $2, $3, $4, $5);

COMMIT;
```
