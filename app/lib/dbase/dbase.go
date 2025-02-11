package dbase

import (
	"context"
	"database/sql"
	"fmt"
	"github/kijunpos/app/lib/apm"
	"github/kijunpos/app/lib/logger"
	"github/kijunpos/config/db"

	"github.com/jmoiron/sqlx"
)

type TableName string

func (tn TableName) String() string {
	return string(tn)
}

// get column name with prefix table name
func (tn TableName) Column(column string) string {
	return tn.String() + "." + column
}

// get column names with prefix table name
func (tn TableName) Columns(columns ...string) []string {
	var result []string
	for _, column := range columns {
		result = append(result, tn.Column(column))
	}
	return result
}

// represent row arrangement query
type RowsArrangement struct {
	Offset  int
	Limit   int
	OrderBy map[string]string // key: column name, value: ASC or DESC
}

type (
	SQLExec interface {
		sqlx.Execer
		sqlx.ExecerContext
		NamedExec(query string, arg interface{}) (sql.Result, error)
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	}

	SQLQuery interface {
		sqlx.Queryer
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	}

	SQLQueryExec interface {
		SQLExec
		SQLQuery
		Rebind(query string) string
	}

	WrapTransactionFunc func(ctx context.Context) error
)

func CreateTrxKey(name db.PSQLName) string {
	return fmt.Sprintf("%s_transaction", name)
}

func GetTxFromContext(ctx context.Context, name db.PSQLName) *sqlx.Tx {
	if tx, ok := ctx.Value(CreateTrxKey(name)).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

func BeginTransaction(ctx context.Context, conn *db.Connection, fn WrapTransactionFunc, isolation ...sql.IsolationLevel) (err error) {
	ctx, span := apm.GetTracer().Start(ctx, "dbase.BeginTransaction")
	defer span.End()

	traceID := span.SpanContext().TraceID()

	// set isolation level, if any
	var il sql.IsolationLevel = sql.LevelReadCommitted
	if len(isolation) == 1 {
		il = isolation[0]
	}

	tx := GetTxFromContext(ctx, conn.Name)
	if tx == nil { // if tx is not exist, then create a new one
		tx, err = conn.DB.BeginTxx(ctx, &sql.TxOptions{
			Isolation: il,
		})
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		defer func() {
			if p := recover(); p != nil || err != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					logger.GetLogger().Errorf("error rollback: %v; trace: %s", errRollback, traceID)
				}

				if p != nil {
					panic(p)
				}
			} else {
				errCommit := tx.Commit()
				if errCommit != nil {
					logger.GetLogger().Errorf("error commit: %v; trace: %s", errCommit, traceID)
				}
			}
		}()

		// set tx
		ctx = context.WithValue(ctx, CreateTrxKey(conn.Name), tx)
	}

	return fn(ctx)
}

func GetExecutor(ctx context.Context, conn *db.Connection) SQLQueryExec {
	var exec SQLQueryExec = conn.DB
	if tx := GetTxFromContext(ctx, conn.Name); tx != nil {
		exec = tx
	}

	return exec
}
