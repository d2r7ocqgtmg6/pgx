// Package pgconn is a low-level PostgreSQL database driver.
//
// It operates at nearly the same level as the C library libpq. It is primarily
// intended to serve as the foundation for higher level libraries such as
// github.com/jackc/pgx. Applications should handle normal queries with a
// higher level library and only use pgconn directly when required for
// low-level access to PostgreSQL functionality.
package pgconn

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// PgError represents an error reported by the PostgreSQL server. See
// http://www.postgresql.org/docs/11/static/protocol-error-fields.html for
// detailed field description.
type PgError struct {
	Severity         string
	Code             string
	Message          string
	Detail           string
	Hint             string
	Position         int32
	InternalPosition int32
	InternalQuery    string
	Where            string
	SchemaName       string
	TableName        string
	ColumnName       string
	DataTypeName     string
	ConstraintName   string
	File             string
	Line             int32
	Routine          string
}

func (pe *PgError) Error() string {
	return pe.Severity + ": " + pe.Message + " (SQLSTATE " + pe.Code + ")"
}

// ConnConfig contains all the options used to establish a connection. It must
// be created by ParseConfig and then it can be modified. A manually
// initialized ConnConfig will cause ConnectConfig to panic.
type ConnConfig struct {
	Host           string
	Port           uint16
	Database       string
	User           string
	Password       string
	TLSConfig      *tls.Config // nil disables TLS
	DialFunc       DialFunc
	ConnectTimeout time.Duration
	RuntimeParams  map[string]string

	createdByParseConfig bool // guard against misconfiguration
}

// DialFunc is a function that can be used to connect to a PostgreSQL server.
type DialFunc func(ctx context.Context, network, addr string) (net.Conn, error)

// PgConn is a low-level PostgreSQL connection handle. It is not safe for
// concurrent usage.
type PgConn struct {
	conn          net.Conn
	config        *ConnConfig
	status        byte // 'I' for idle, 'T' for in transaction, 'E' for in failed transaction
	pid           uint32
	secretKey     uint32
	parameterStatuses map[string]string
	txStatus      byte
	closed        bool
}

// Connect establishes a connection to a PostgreSQL server using the
// environment and connString to provide configuration. See documentation for
// ParseConfig for details. ctx can be used to cancel a connect attempt.
func Connect(ctx context.Context, connString string) (*PgConn, error) {
	config, err := ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	return ConnectConfig(ctx, config)
}

// ConnectConfig establishes a connection to a PostgreSQL server using config.
// config must have been created by ParseConfig. ctx can be used to cancel a
// connect attempt.
func ConnectConfig(ctx context.Context, config *ConnConfig) (*PgConn, error) {
	if !config.createdByParseConfig {
		panic("config must be created by ParseConfig")
	}

	network := "tcp"
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	var dialFunc DialFunc
	if config.DialFunc != nil {
		dialFunc = config.DialFunc
	} else {
		dialer := &net.Dialer{
			Timeout: config.ConnectTimeout,
		}
		dialFunc = dialer.DialContext
	}

	netConn, err := dialFunc(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	pgConn := &PgConn{
		conn:              netConn,
		config:            config,
		parameterStatuses: make(map[string]string),
	}

	return pgConn, nil
}

// Close closes a connection. It is safe to call Close on a already closed
// connection.
func (c *PgConn) Close(ctx context.Context) error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.conn.Close()
}

// IsClosed reports if the connection has been closed.
func (c *PgConn) IsClosed() bool {
	return c.closed
}

// ParameterStatus returns the value of a parameter reported by the server
// (e.g. server_version). Returns an empty string for unknown parameters.
func (c *PgConn) ParameterStatus(key string) string {
	return c.parameterStatuses[key]
}
