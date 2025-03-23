package db

import (
	"fmt"
	"go.opentelemetry.io/otel/sdk/trace"
	"reflect"
	"unsafe"

	"github.com/NawafSwe/media-scout-service/cmd/config"
	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// NewDBConn creates a new db connection and returns *sqlx.DB.
func NewDBConn(cfg config.DB, tracer *trace.TracerProvider) (*sqlx.DB, error) {
	var opts []otelsql.Option
	if tracer != nil {
		opts = append(opts, otelsql.WithTracerProvider(tracer), otelsql.WithSpanOptions(otelsql.SpanOptions{
			OmitConnResetSession: true,
		}))
	}
	db, err := otelsql.Open("postgres", cfg.DSN, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a db conn: %w", err)
	}
	if cfg.MaxOpenConnections != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}
	if cfg.MaxIdleConnections != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}
	if cfg.MaxConnectionsLifetime != 0 {
		db.SetConnMaxLifetime(cfg.MaxConnectionsLifetime)
	}
	// Create an instrumented driver
	drvInstrum := otelsql.WrapDriver(db.Driver())

	// Extract accessible/writable connector since it's not exported and cannot be accessed directly
	connFld := reflect.ValueOf(db).Elem().FieldByName("connector")
	connFldMod := reflect.NewAt(connFld.Type(), unsafe.Pointer(connFld.UnsafeAddr())).Elem()

	// Make an addressable value of the connector
	connAddr := reflect.New(reflect.ValueOf(connFldMod.Interface()).Type()).Elem()
	connAddr.Set(reflect.ValueOf(connFldMod.Interface()))

	// Extract driver value from connector and make it writeable
	drvFld := connAddr.FieldByName("driver")
	drvFldAcc := reflect.NewAt(drvFld.Type(), unsafe.Pointer(drvFld.UnsafeAddr())).Elem()

	// Replace driver value with instrumented one
	drvFldAcc.Set(reflect.ValueOf(drvInstrum))
	// Replace connector value with new instrumented connector
	connFldMod.Set(connAddr)

	return sqlx.NewDb(db, "postgres"), nil
}
