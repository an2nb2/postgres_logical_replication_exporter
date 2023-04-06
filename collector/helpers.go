package collector

import "github.com/jackc/pgx"

func calculateLag(lsn1, lsn2 string) (uint64, error) {
	val1, err := pgx.ParseLSN(lsn1)
	if err != nil {
		return 0, err
	}
	val2, err := pgx.ParseLSN(lsn2)
	if err != nil {
		return 0, err
	}
	return val1 - val2, nil
}

func getStatus(active bool) (st uint8) {
	if active {
		st = 1
	}
	return
}
