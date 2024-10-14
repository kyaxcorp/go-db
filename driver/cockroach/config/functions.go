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
	host := c.masterConfig.Credentials.Host
	user := c.masterConfig.Credentials.User
	password := c.masterConfig.Credentials.Password
	dbName := c.masterConfig.Credentials.DbName
	port := c.masterConfig.Credentials.Port

	timeZone := c.masterConfig.Credentials.TimeZone
	sslMode := c.masterConfig.Credentials.SSLMode

	// TODO: continue SSL
	/*sslFactory := c.masterConfig.Credentials.SSLFactory
	caCertificate := c.masterConfig.Credentials.CACertificate
	clientCertificate := c.masterConfig.Credentials.ClientCertificate
	clientPrivateKey := c.masterConfig.Credentials.ClientPrivateKey*/

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

	if conv.ParseBool(co.TimeZone) {
		timeZone = c.TimeZone
	}

	if conv.ParseBool(co.SSLMode) {
		sslMode = c.SSLMode
	}

	// TODO: continue SSL

	/*sslFactory = c.SSLFactory
	caCertificate = c.CACertificate
	clientCertificate = c.ClientCertificate
	clientPrivateKey = c.ClientPrivateKey*/

	//log.Println("parse time", c.masterConfig.ParseTime)

	dsn1 := "host=" + host +
		" user=" + user +
		" password="

	dsn2 := " dbname=" + dbName +
		" port=" + strconv.Itoa(port) +
		" sslmode=" + sslMode

	// if not empty... then add the timezone
	if timeZone != "" {
		dsn2 = dsn2 + " TimeZone=" + timeZone
	}

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
