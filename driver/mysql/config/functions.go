package config

import (
	"strconv"
	"time"

	"github.com/kyaxcorp/go-helper/conv"
	"github.com/kyaxcorp/go-helper/hash"
)

type DSN struct {
	Plain       string
	Secured     string
	Checksum    string
	GeneratedAt time.Time
}

func (c *Connection) GenerateDSN() DSN {
	charset := c.masterConfig.Charset
	if charset == "" {
		// utf8mb4 is 4 byte char & utf8 is 3 byte char.
		// the utf8mb4 is newer and can store any unicode char!
		charset = "utf8mb4"
	}

	host := c.masterConfig.Credentials.Host
	user := c.masterConfig.Credentials.User
	password := c.masterConfig.Credentials.Password
	dbName := c.masterConfig.Credentials.DbName
	port := c.masterConfig.Credentials.Port

	co := c.CredentialsOverrides

	if conv.ParseBool(co.Host) {
		host = c.Host
	}

	if conv.ParseBool(co.User) {
		user = c.User
	}

	if conv.ParseBool(co.Password) {
		password = c.Password
	}

	if conv.ParseBool(co.DbName) {
		dbName = c.DbName
	}

	if conv.ParseBool(co.Port) {
		port = c.Port
	}

	//log.Println("parse time", c.masterConfig.ParseTime)

	dsn1 := user + ":"

	dsn2 := "@tcp(" +
		host + ":" +
		strconv.Itoa(port) + ")/" +
		dbName + "?charset=" +
		charset + "&parseTime=" +
		strconv.FormatBool(conv.ParseBool(c.masterConfig.ParseTime)) +
		"&loc=Local"

	dsnPlain := dsn1 + password + dsn2
	dsnSecured := dsn1 + "*********" + dsn2
	dsnHash := hash.MD5(dsnPlain)

	return DSN{
		Plain:       dsnPlain,
		Secured:     dsnSecured,
		Checksum:    dsnHash,
		GeneratedAt: time.Now(),
	}
}
