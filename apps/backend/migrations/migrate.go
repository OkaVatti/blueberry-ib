// apps/backend/migrations/migrate.go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		filename TEXT NOT NULL UNIQUE,
		applied_at TIMESTAMP NOT NULL
	);
	`)
	return err
}

func alreadyApplied(db *sql.DB, filename string) (bool, error) {
	var cnt int
	err := db.QueryRow("SELECT COUNT(1) FROM schema_migrations WHERE filename = $1", filename).Scan(&cnt)
	return cnt > 0, err
}

func recordApplied(db *sql.DB, filename string) error {
	_, err := db.Exec("INSERT INTO schema_migrations (filename, applied_at) VALUES ($1, $2)", filename, time.Now().UTC())
	return err
}

func main() {
	dsn := flag.String("dsn", "", "Postgres DSN (e.g. postgres://user:pass@host:5432/db?sslmode=disable)")
	dir := flag.String("dir", "./migrations", "migrations directory (contains .sql files)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("dsn is required (use --dsn)")
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("ensure migrations table: %v", err)
	}

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatalf("read dir: %v", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, fname := range sqlFiles {
		applied, err := alreadyApplied(db, fname)
		if err != nil {
			log.Fatalf("checking applied: %v", err)
		}
		if applied {
			fmt.Printf("skipping already applied: %s\n", fname)
			continue
		}

		raw, err := ioutil.ReadFile(filepath.Join(*dir, fname))
		if err != nil {
			log.Fatalf("read file: %v", err)
		}
		// split on semicolon; allow multiple statements
		parts := strings.Split(string(raw), ";")
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("begin tx: %v", err)
		}
		for _, p := range parts {
			sqlstmt := strings.TrimSpace(p)
			if sqlstmt == "" {
				continue
			}
			if _, err := tx.Exec(sqlstmt); err != nil {
				_ = tx.Rollback()
				log.Fatalf("exec statement in %s: %v\nSQL: %s", fname, err, sqlstmt)
			}
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("commit: %v", err)
		}
		if err := recordApplied(db, fname); err != nil {
			log.Fatalf("record applied: %v", err)
		}
		fmt.Printf("applied migration: %s\n", fname)
	}
	fmt.Println("migrations done")
}
