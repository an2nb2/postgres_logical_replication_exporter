package pg

import (
	"database/sql"
	"errors"
)

type ReplicationSlot struct {
	Name              string `db:"slot_name"`
	RestartLsn        string `db:"restart_lsn"`
	ConfirmedFlushLsn string `db:"confirmed_flush_lsn"`
	Active            bool   `db:"active"`
}

type Publication struct {
	Active  bool          `db:"active"`
	Name    string        `db:"application_name"`
	Pid     sql.NullInt64 `db:"pid"`
	SentLsn string        `db:"sent_lsn"`
	State   string        `db:"state"`
	Tmp     bool          `db:"temporary"`
}

// Returns a list of postgres publications that already have replication slot created, i.e.
// there is a subscription on standby instance side.
func (db *DB) Publications() ([]Publication, error) {
	var pubs []Publication
	query := `
  SELECT active, application_name, pid, sent_lsn, state, temporary
  FROM pg_stat_replication
  JOIN pg_replication_slots ON pg_stat_replication.pid = pg_replication_slots.active_pid
  `
	err := db.Select(&pubs, query)
	if err == sql.ErrNoRows {
		return pubs, nil
	}

	return pubs, err
}

// Returns the current write-ahead log write location as a string.
// When no replication exists it will return an error.
func (db *DB) CurrentWalLsn() (string, error) {
	var val sql.NullString
	query := `
  SELECT pg_current_wal_lsn()
  `
	err := db.Get(&val, query)
	if !val.Valid {
		err = errors.New("current wal lsn is null")
	}

	return val.String, err
}

// Returns a list of postgres replication slots.
func (db *DB) ReplicationSlots() ([]ReplicationSlot, error) {
	var slots []ReplicationSlot
	query := `
  SELECT slot_name, restart_lsn, confirmed_flush_lsn, active FROM pg_replication_slots
  `
	err := db.Select(&slots, query)
	if err == sql.ErrNoRows {
		return slots, nil
	}

	return slots, err
}
