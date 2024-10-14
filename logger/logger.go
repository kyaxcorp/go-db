package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-helper/conv"
	"github.com/kyaxcorp/go-helper/errors2"
	"github.com/kyaxcorp/go-logger"
	loggerConfig "github.com/kyaxcorp/go-logger/config"
	"github.com/kyaxcorp/go-logger/model"
	"github.com/kyaxcorp/go-logger/paths"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLogger struct {
	Logger       *model.Logger
	dbType       string
	instanceName string
}

func NewLogger(
	instanceName string,
	dbType string,
	logConfig loggerConfig.Config,
) *GormLogger {
	// Create the logger for the database or it should be the same?!
	/*logConfig, _ := loggerConfig.DefaultConfig(&loggerConfig.Config{
		IsEnabled:  "yes",
		Name:       dbType,
		ModuleName: "database-modulename",
		DirLogPath: paths.GetDatabasePath(dbType),
	})*/

	logCfg, _ := loggerConfig.DefaultConfig(&logConfig)
	// If there is no Name, we set the DB User and DB Name
	if logCfg.Name == "" {
		logCfg.Name = instanceName
	}

	// If there is no specific path where to save
	if logCfg.DirLogPath == "" {
		logCfg.DirLogPath = paths.GetDatabasePath(dbType + "/" + instanceName)
	}

	// TODO: maybe we should set this as a separate string?!!.. added below in the logs...
	if logCfg.ModuleName == "" {
		logCfg.ModuleName = dbType
	}

	return &GormLogger{
		dbType:       dbType,
		instanceName: instanceName,
		Logger:       logger.New(logCfg),
	}
}

func (l *GormLogger) l() *zerolog.Logger {
	return l.Logger.Logger
}

func (l *GormLogger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	l.print(l.l().Info(), ctx, msg, data...)
}

func (l *GormLogger) Warn(
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	l.print(l.l().Warn(), ctx, msg, data...)
}

func (l *GormLogger) print(
	event *zerolog.Event,
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	event.
		Str("db_instance", l.instanceName).
		Str("file", utils.FileWithLineNum()).
		Msg(fmt.Sprintf(msg, data...))
}

func (l *GormLogger) Error(
	ctx context.Context,
	msg string,
	data ...interface{},
) {
	// We will print all gorm errors... even ErrRecordNotFound
	// if you don't need to see these errors, then you can simply set the log level to another level or even disable it!

	l.print(
		l.l().Error(),
		ctx,
		msg,
		data...,
	)

	errors2.NewCustom(errors2.CustomError{
		Code:    0,
		Message: "Database Error -> " + l.dbType + " -> " + l.instanceName + " -> " + msg,
		// We don't need to log this to app logs... because the logger does already
		LogToApp: false,
	})
}

func (l *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	_err error,
) {
	elapsed := time.Since(begin).Milliseconds()
	elapsedStr := conv.Int64ToStr(elapsed)
	sql, rows := fc()
	sqlFile := utils.FileWithLineNum()

	elapsedInfo := color.Style{color.LightRed}.Render("sql:" + elapsedStr + "ms")

	// db_instance has being added for Application specific logs

	if _err != nil && !errors.Is(_err, gorm.ErrRecordNotFound) {
		l.l().Error().
			Str("db_instance", l.instanceName).
			Str("sql_file", sqlFile).
			Str("sql", sql).
			Int64("rows", rows).
			Time("sql_begin_time", begin).
			Int64("sql_elapsed_time_ms", elapsed).
			Err(_err).Msg("" + elapsedInfo)
		// Sentry
		errors2.NewCustom(errors2.CustomError{
			Code: 0,
			Message: "Database Error -> " +
				"db type -> " + l.dbType + " -> " +
				"db instance -> " + l.instanceName + " -> " +
				"sql_file -> " + sqlFile + " -> " +
				"sql -> " + sql + " -> " +
				"rows -> " + conv.Int64ToStr(rows) + " -> " +
				"sql_elapsed_time_ms -> " + conv.Int64ToStr(elapsed) + " -> " +
				"sql_begin_time -> " + begin.String() + " -> " +
				"error -> " + _err.Error(),
			// We don't need to log this to app logs... because the logger does already
			LogToApp: false,
		})
	} else {
		// It's simply a trace...
		l.l().Info().
			Str("db_instance", l.instanceName).
			Str("sql_file", sqlFile).
			Str("sql", sql).
			Int64("rows", rows).
			Time("sql_begin_time", begin).
			Int64("sql_elapsed_time_ms", elapsed).Msg("" + elapsedInfo)
	}

}
