package dbresolver

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"strings"
)

func (dr *DBResolver) registerCallbacks(db *gorm.DB) {
	// When creating, it should switch the source...

	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.registerCallbacks")
	}
	info().Msg("registering callbacks")
	defer info().Msg("leaving...")

	// Switching sources
	dr.Callback().Create().Before("*").Register("gorm:db_resolver", dr.switchSource)
	dr.Callback().Update().Before("*").Register("gorm:db_resolver", dr.switchSource)
	dr.Callback().Delete().Before("*").Register("gorm:db_resolver", dr.switchSource)

	// Switching replicas
	dr.Callback().Query().Before("*").Register("gorm:db_resolver", dr.switchReplica)
	dr.Callback().Row().Before("*").Register("gorm:db_resolver", dr.switchReplica)

	// Switching on guess
	dr.Callback().Raw().Before("*").Register("gorm:db_resolver", dr.switchGuess)
}

func (dr *DBResolver) switchSource(db *gorm.DB) {
	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.switchSource")
	}
	info().Msg("switching source")
	defer info().Msg("leaving...")

	/*dbc, _err := db.Begin().DB()


	v := reflect.ValueOf(dbc)
	dsnReflection := v.FieldByName("DSN")
	dsn := dsnReflection.String()

	if dsn, ok := db.Dialector.(interface{ DSN }); ok {

	}*/
	/*db.Transaction(func(tx *gorm.DB) error {
		tx.Select("")
		sqld,_err := tx.DB()

	})*/
	if !isTransaction(db.Statement.ConnPool) {
		info().Msg("it's not a transaction, so executing...")
		db.Statement.ConnPool = dr.resolve(db.Statement, Write)
	}
}

func (dr *DBResolver) switchReplica(db *gorm.DB) {
	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.switchReplica")
	}
	info().Msg("switching replica")
	defer info().Msg("leaving...")

	if !isTransaction(db.Statement.ConnPool) {
		info().Msg("it's not a transaction, so continuing...")
		if rawSQL := db.Statement.SQL.String(); len(rawSQL) > 0 {
			dr.switchGuess(db)
		} else {
			_, locking := db.Statement.Clauses["FOR"]
			if _, ok := db.Statement.Clauses[writeName]; ok || locking {
				db.Statement.ConnPool = dr.resolve(db.Statement, Write)
			} else {
				db.Statement.ConnPool = dr.resolve(db.Statement, Read)
			}
		}
	}
}

func (dr *DBResolver) switchGuess(db *gorm.DB) {
	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.switchGuess")
	}
	info().Msg("switch guess...")
	defer info().Msg("leaving...")

	if !isTransaction(db.Statement.ConnPool) {
		info().Msg("it's not a transaction, continuing...")
		if _, ok := db.Statement.Clauses[writeName]; ok {
			db.Statement.ConnPool = dr.resolve(db.Statement, Write)
		} else if rawSQL := strings.TrimSpace(db.Statement.SQL.String()); len(rawSQL) > 10 && strings.EqualFold(rawSQL[:6], "select") && !strings.EqualFold(rawSQL[len(rawSQL)-10:], "for update") {
			db.Statement.ConnPool = dr.resolve(db.Statement, Read)
		} else {
			db.Statement.ConnPool = dr.resolve(db.Statement, Write)
		}
	}
}

func isTransaction(connPool gorm.ConnPool) bool {
	_, ok := connPool.(gorm.TxCommitter)
	return ok
}
