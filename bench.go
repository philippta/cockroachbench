package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

const count = 1_000

func main() {
	connCrdb := connectCockroach()
	runtest(connCrdb, "COCKROACH")
	connCrdb.Close(context.Background())

	fmt.Println()

	connPsql := connectPostgres()
	runtest(connPsql, "POSTGRES")
	connPsql.Close(context.Background())
}

func runtest(conn *pgx.Conn, label string) {
	fmt.Println("====================================")
	fmt.Println(label)
	fmt.Println("------------------------------------")
	droptable(conn)
	createtable(conn)
	insert(conn)
	fmt.Println("------------------------------------")
	randomreads(conn)
	fmt.Println("====================================")
}

func droptable(conn *pgx.Conn) {
	_, err := conn.Exec(context.TODO(), `drop table if exists foo`)
	must(err)
}

func createtable(conn *pgx.Conn) {
	_, err := conn.Exec(context.TODO(), `create table if not exists foo(a integer primary key, b integer, c integer)`)
	must(err)
}

func insert(conn *pgx.Conn) {
	fmt.Println("Benchmarking inserts...")

	now := time.Now()
	for i := 0; i < count; i++ {
		_, err := conn.Exec(context.Background(), "insert into foo values ($1, $2, $3)", i, i*2, i*3)
		must(err)
	}
	dur := time.Since(now)
	fmt.Printf("Inserts:     %v\n", count)
	fmt.Printf("Duration:    %v\n", dur)
	fmt.Printf("Inserts/Sec: %v\n", float64(count)/float64(dur.Seconds()))
}

func randomreads(conn *pgx.Conn) {
	fmt.Println("Benchmarking random reads...")

	now := time.Now()
	for i := 0; i < count; i++ {
		var a, b, c int64
		err := conn.QueryRow(context.Background(), "select * from foo where a = $1", rand.Intn(count)).Scan(&a, &b, &c)
		must(err)
	}
	dur := time.Since(now)
	fmt.Printf("Reads:     %v\n", count)
	fmt.Printf("Duration:  %v\n", dur)
	fmt.Printf("Reads/Sec: %v\n", float64(count)/float64(dur.Seconds()))
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
