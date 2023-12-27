package pg

import (
	"database/sql"
)

type Subscription struct {
	Pid                sql.NullInt64  `db:"pid"`
	Name               string         `db:"subname"`
	Relname            sql.NullString `db:"relname"`
	ReceivedLsn        string         `db:"received_lsn"`
	LastMsgReceiptTime string         `db:"last_msg_receipt_time"`
}

// Returns a list of postgres subscriptions.
func (db *DB) Subscriptions() ([]Subscription, error) {
	var subs []Subscription
	query := `
  SELECT pid, subname, received_lsn, last_msg_receipt_time, relname
  FROM pg_stat_subscription
  LEFT JOIN pg_stat_all_tables on pg_stat_subscription.relid = pg_stat_all_tables.relid
  `
	err := db.Select(&subs, query)
	if err == sql.ErrNoRows {
		return subs, nil
	}
	return subs, err
}
