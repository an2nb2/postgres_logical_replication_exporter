package collector

import "github.com/jackc/pgx"

// Takes two string arguments in format of LSN (Log Sequence Number) and returns the difference
// between them. Returns an error if it fails to parse one of the arguments.
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

// Takes boolean argument and converts it to uint8 which is used by prometheus metrics.
func getStatus(active bool) (st uint8) {
	if active {
		st = 1
	}
	return
}
