package inhire

import (
	"context"
	"database/sql"
	"fmt"
)

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

type Store struct {
	db *sql.DB
}

func (store *Store) CachedJobs(ctx context.Context) ([]JobInfo, error) {
	sql := "SELECT position, url, page_id FROM jobs;"
	rows, err := store.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]JobInfo, 0)
	for rows.Next() {
		var position, url, pageid string
		err := rows.Scan(&position, &url, &pageid)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, JobInfo{
			ID:           pageid,
			PageURL:      url,
			PositionName: position,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (store *Store) SaveJobs(ctx context.Context, jobs ...JobInfo) error {
	sql := `INSERT INTO jobs (position, page_id, url, checked_at)
					VALUES (?, ?, ?, CURRENT_TIMESTAMP)
					ON CONFLICT(page_id) DO UPDATE SET
					checked_at = CURRENT_TIMESTAMP;`

	for _, j := range jobs {
		_, err := store.db.ExecContext(ctx, sql, j.PositionName, j.ID, j.PageURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) SaveLinks(ctx context.Context, links ...string) error {
	sql := "INSERT OR IGNORE INTO links (url) VALUES (?);"
	for _, l := range links {
		_, err := store.db.ExecContext(ctx, sql, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *Store) CachedLinks(ctx context.Context) ([]string, error) {
	sql := "SELECT url FROM links;"
	rows, err := store.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []string
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		links = append(links, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println(len(links))
	return links, nil
}
