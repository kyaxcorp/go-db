package dberr

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsDuplicateKey(_err error) bool {
	if _err == nil {
		return false
	}
	switch t := _err.(type) {
	case *pgconn.PgError:
		switch t.Code {
		case "23505":
			return true
		}
	case *mysql.MySQLError:
		switch t.Number {
		case 1062:
			return true
		}
	}
	return false
}
