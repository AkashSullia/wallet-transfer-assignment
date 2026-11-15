package database

import (
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

func GetPGError(err error) (error, bool) {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case UniqueViolation:
			re := regexp.MustCompile(`(?m)\=\((.*)\)`)
			entry := string(re.Find([]byte(pgErr.Detail)))
			return fmt.Errorf("unique key violation for entry %s", entry), true
		case ForeignKeyViolation:
			return fmt.Errorf("foreign key violation %s", pgErr.Message), true
		}
	}
	return nil, false
}

func IsUniqueKeyViolation(err error) (bool, string) {
	pgErr, ok := err.(*pgconn.PgError)
	re := regexp.MustCompile(`\=\((.*)\)`)
	key := re.FindStringSubmatch(pgErr.Detail)
	return ok && (pgErr.Code == UniqueViolation), key[1]
}
