package pg

import (
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
	Hostname string
	Database string
}

// Tries to establish connection with database using given Connection URI and
// returns a wrapped sqlx.DB structure with few more fields.
// If it doesn't succeed to connect then it causes to exit with code 1.
func MustConnect(uri string, poolSize int) *DB {
	config, err := pgx.ParseURI(uri)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse postgres uri string: %v\n", err)
		os.Exit(1)
	}

	pgxdb := stdlib.OpenDB(config)
	pgxdb.SetMaxOpenConns(poolSize)
	pgxdb.SetConnMaxLifetime(time.Duration(5) * time.Second)

	if err := pgxdb.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &DB{
		sqlx.NewDb(pgxdb, "pgx"),
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		config.Database,
	}
}
