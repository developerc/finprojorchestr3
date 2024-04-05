package sqlite

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

//var db *sql.DB

func CreateSqliteDb() error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		//panic(err)
		return errors.New("can't open db")
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		//panic(err)
		return errors.New("can't ping context")
	}
	if err = createTables(ctx, db); err != nil {
		//panic(err)
		return errors.New("can't create tables")
	}
	return nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	const (
		tasksTable = `
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER,
		agentid INTEGER,
		status TEXT,
		expr TEXT,
		result FLOAT,
		begindate timestamp,
		enddate timestamp
	);`

		agentsTable = `
	CREATE TABLE IF NOT EXISTS agents(
		id INTEGER,
		ip TEXT,
		port INTEGER
	);`
	)

	if _, err := db.ExecContext(ctx, tasksTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, agentsTable); err != nil {
		return err
	}
	return nil
}
