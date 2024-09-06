package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pretodev/inhireapp/config/env"
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

func OpenConn(ctx context.Context, cfg env.Config) (*sql.DB, error) {
	source := fmt.Sprintf("file:%s?cache=shared&mode=rwc", cfg.SQLiteDBPath)
	db, err := sql.Open("sqlite3", source)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	_, err = db.ExecContext(ctx, schema)
	return db, err
}
