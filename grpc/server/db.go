package server

import (
	"context"
	"database/sql"
	"errors"

	pb "github.com/developerc/finprojorchestr3/proto"
	_ "github.com/mattn/go-sqlite3"
)

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

func InsertTask(task pb.Task) error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		//panic(err)
		return errors.New("can't open db")
	}
	defer db.Close()

	var q = `
	INSERT INTO tasks (id, agentid, status, expr, result) values ($1, $2, $3, $4, $5)
	`
	_, err = db.ExecContext(ctx, q, task.Id, task.Agentid, task.Status, task.Expr, task.Result)
	if err != nil {
		return err
	}
	return nil
}
