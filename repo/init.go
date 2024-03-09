package repo

import (
	"context"
	"fmt"
	"log"
	"time"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitemigration"
)

type Repo struct {
	pool    *sqlitemigration.Pool
	modules map[string]*Module
}

func New() (*Repo, error) {
	// Open a database fil, db.sqlite3:
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := poolMigration(ctx, "file:db.sqlite3")
	if err != nil {
		return nil, fmt.Errorf("sqlite.OpenConn: %w", err)
	}

	r := &Repo{
		pool:    conn,
		modules: make(map[string]*Module),
	}
	return r, nil
}

func getMigrations() []string {
	return []string{
		`CREATE TABLE modules ( id INTEGER NOT NULL PRIMARY KEY, name VARCHAR(255) );`,
		`CREATE TABLE packages ( id INTEGER NOT NULL PRIMARY KEY, name VARCHAR(255));`,
	}
}

func poolMigration(ctx context.Context, path string) (*sqlitemigration.Pool, error) {
	schema := sqlitemigration.Schema{
		// Each element of the Migrations slice is applied in sequence. When you
		// want to change the schema, add a new SQL script to this list.
		//
		// Existing databases will pick up at the same position in the Migrations
		// slice as they last left off.
		Migrations: getMigrations(),

		// The RepeatableMigration is run after all other Migrations if any
		// migration was run. It is useful for creating triggers and views.
		RepeatableMigration: "",
	}
	// Open a pool. This does not block, and will start running any migrations
	// asynchronously.
	pool := sqlitemigration.NewPool(path, schema, sqlitemigration.Options{
		Flags: sqlite.OpenReadWrite | sqlite.OpenCreate,
		OnError: func(e error) {
			log.Println(e)
		},
	})
	// Get a connection. This blocks until the migration completes.
	conn, err := pool.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Get: %w", err)
	}
	defer pool.Put(conn)

	return pool, nil
}
