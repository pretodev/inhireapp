package db

import (
	"context"
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

const schema = `
CREATE TABLE IF NOT EXISTS links(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  url TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jobs(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  position TEXT NOT NULL,
  page_id TEXT UNIQUE NOT NULL,
  url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  checked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func Connect(ctx context.Context) *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", "file:example.db?cache=shared&mode=rwc")
		if err != nil {
			log.Fatalf("failed to open database: %v", err)
		}
		_, err = db.ExecContext(ctx, schema)
		if err != nil {
			log.Fatalf("failed to create schema: %v", err)
		}
	})

	return db
}
