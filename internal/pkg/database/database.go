package database

import (
	"dbms/internal/pkg/utils/log"

	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBbyterow [][]byte

type ConnectionPool interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type DBManager struct {
	Pool ConnectionPool
}

func InitDatabase() *DBManager {
	return &DBManager{
		Pool: nil,
	}
}

func (dbm *DBManager) Connect() {
	var connString string = "postgres://docker:docker@127.0.0.1:5432/docker?sslmode=disable&pool_max_conns=1000"

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Warn("{Connect} Postgres error")
		log.Error(err)
		return
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Warn("{Connect} Ping error")
		log.Error(err)
		return
	}

	log.Info("Successful connection to postgres")
	log.Info("Connection params: " + connString)
	dbm.Pool = pool
}

func (dbm *DBManager) Disconnect() {
	dbm.Pool.Close()
	log.Info("Postgres disconnected")
}

func (dbm *DBManager) Query(queryString string, params ...interface{}) ([]DBbyterow, error) {
	transactionContext := context.Background()
	tx, err := dbm.Pool.Begin(transactionContext)
	if err != nil {
		log.Warn("{Query} Error connecting to a pool")
		log.Error(err)
		return nil, err
	}

	defer func() {
		err := tx.Rollback(transactionContext)
		if err != nil {
			log.Error(err)
		}
	}()

	rows, err := tx.Query(transactionContext, queryString, params...)
	if err != nil {
		log.Warn("{Query} Error in query: " + queryString)
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([]DBbyterow, 0)
	for rows.Next() {
		rowBuffer := make(DBbyterow, 0)
		rowBuffer = append(rowBuffer, rows.RawValues()...)
		result = append(result, rowBuffer)
	}

	err = tx.Commit(transactionContext)
	if err != nil {
		log.Warn("{Query} Error committing")
		log.Error(err)
		return nil, err
	}

	return result, nil
}
