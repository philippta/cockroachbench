package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	count       = 5000
	numAccounts = 5000
	numBranches = 500
	numTellers  = 500
)

func main() {
	connPsql := connectPostgres()
	runtest(connPsql, "POSTGRES")
	connPsql.Close(context.Background())

	fmt.Println()

	connCrdbSingle := connectCockroachSingleNode()
	runtest(connCrdbSingle, "COCKROACH SINGLE NODE")
	connCrdbSingle.Close(context.Background())

	fmt.Println()

	connCrdb := connectCockroach()
	runtest(connCrdb, "COCKROACH 3 NODES")
	connCrdb.Close(context.Background())

	fmt.Println()
}

func runtest(conn *pgx.Conn, label string) {
	fmt.Println("====================================")
	fmt.Println(label)
	fmt.Println("------------------------------------")
	droptables(conn)
	createtables(conn)
	seedTables(conn, numAccounts, numBranches, numTellers)
	runTransactions(conn, count, numAccounts, numBranches, numTellers)
	fmt.Println("====================================")
}

func droptables(conn *pgx.Conn) {
	queries := []string{
		"DROP TABLE IF EXISTS pgbench_accounts;",
		"DROP TABLE IF EXISTS pgbench_branches;",
		"DROP TABLE IF EXISTS pgbench_tellers;",
		"DROP TABLE IF EXISTS pgbench_history;",
	}

	for _, query := range queries {
		_, err := conn.Exec(context.TODO(), query)
		must(err)
	}
}

func createtables(conn *pgx.Conn) {
	queries := []string{
		`CREATE TABLE pgbench_accounts (
			aid        INTEGER NOT NULL,
			bid        INTEGER NOT NULL,
			abalance   INTEGER NOT NULL,
			filler     CHAR(84) NOT NULL
		);`,
		`CREATE TABLE pgbench_branches (
			bid        INTEGER NOT NULL,
			bbalance   INTEGER NOT NULL,
			filler     CHAR(88) NOT NULL
		);`,
		`CREATE TABLE pgbench_tellers (
			tid        INTEGER NOT NULL,
			bid        INTEGER NOT NULL,
			tbalance   INTEGER NOT NULL,
			filler     CHAR(84) NOT NULL
		);`,
		`CREATE TABLE pgbench_history (
			tid        INTEGER NOT NULL,
			bid        INTEGER NOT NULL,
			aid        INTEGER NOT NULL,
			delta      INTEGER NOT NULL,
			mtime      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		_, err := conn.Exec(context.TODO(), query)
		must(err)
	}
}

func executeTransaction(conn *pgx.Conn, aid, tid, bid, delta int) {
	tx, err := conn.Begin(context.TODO()) // Start a transaction
	must(err)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(context.TODO())
			panic(p)
		} else {
			err = tx.Commit(context.TODO())
			must(err)
		}
	}()

	// Update the account balance
	_, err = tx.Exec(context.TODO(), `
		UPDATE pgbench_accounts
		SET abalance = abalance + $1
		WHERE aid = $2;
	`, delta, aid)
	must(err)

	// Update the teller balance
	_, err = tx.Exec(context.TODO(), `
		UPDATE pgbench_tellers
		SET tbalance = tbalance + $1
		WHERE tid = $2;
	`, delta, tid)
	must(err)

	// Update the branch balance
	_, err = tx.Exec(context.TODO(), `
		UPDATE pgbench_branches
		SET bbalance = bbalance + $1
		WHERE bid = $2;
	`, delta, bid)
	must(err)

	// Insert a record into the history table
	_, err = tx.Exec(context.TODO(), `
		INSERT INTO pgbench_history (tid, bid, aid, delta, mtime)
		VALUES ($1, $2, $3, $4, $5);
	`, tid, bid, aid, delta, time.Now())
	must(err)

	// log.Printf("Transaction completed: aid=%d, tid=%d, bid=%d, delta=%d", aid, tid, bid, delta)
}

func runTransactions(conn *pgx.Conn, numTransactions, numAccounts, numBranches, numTellers int) {
	fmt.Println("Benchmarking TPC-B...")
	start := time.Now()
	for i := 1; i <= numTransactions; i++ {
		aid := rand.Intn(numAccounts) // Example account ID
		bid := rand.Intn(numBranches) // Example branch ID
		tid := rand.Intn(numTellers)  // Example teller ID
		delta := 100 - i              // Example transaction amount

		executeTransaction(conn, aid, tid, bid, delta)
	}
	dur := time.Since(start)
	fmt.Printf("Transactions:     %v\n", numTransactions)
	fmt.Printf("Duration:         %v\n", dur)
	fmt.Printf("Transactions/Sec: %v\n", float64(numTransactions)/float64(dur.Seconds()))
}

func seedTables(conn *pgx.Conn, numAccounts, numBranches, numTellers int) {
	ctx := context.TODO()

	// Seed pgbench_branches
	for bid := 1; bid <= numBranches; bid++ {
		_, err := conn.Exec(ctx, `
			INSERT INTO pgbench_branches (bid, bbalance, filler)
			VALUES ($1, $2, $3);
		`, bid, rand.Intn(10000), "branch filler data")
		must(err)
	}

	// Seed pgbench_tellers
	for tid := 1; tid <= numTellers; tid++ {
		bid := (tid % numBranches) + 1 // Assign teller to a branch
		_, err := conn.Exec(ctx, `
			INSERT INTO pgbench_tellers (tid, bid, tbalance, filler)
			VALUES ($1, $2, $3, $4);
		`, tid, bid, rand.Intn(1000), "teller filler data")
		must(err)
	}

	// Seed pgbench_accounts
	for aid := 1; aid <= numAccounts; aid++ {
		bid := (aid % numBranches) + 1 // Assign account to a branch
		_, err := conn.Exec(ctx, `
			INSERT INTO pgbench_accounts (aid, bid, abalance, filler)
			VALUES ($1, $2, $3, $4);
		`, aid, bid, rand.Intn(5000), "account filler data")
		must(err)
	}
	fmt.Printf("Seeded tables with:\n- %d accounts\n- %d branches\n- %d tellers\n", numAccounts, numBranches, numTellers)
}

func connectPostgres() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:15432/postgres")
	must(err)
	return conn
}

func connectCockroach() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://root@localhost:26257/defaultdb")
	must(err)
	return conn
}

func connectCockroachSingleNode() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://root@localhost:26256/defaultdb")
	must(err)
	return conn
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
