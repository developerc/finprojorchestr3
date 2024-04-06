package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pb "github.com/developerc/finprojorchestr3/proto"
	_ "github.com/mattn/go-sqlite3"
)

func CreateSqliteDb() error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return errors.New("can't open db")
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		return errors.New("can't ping context")
	}
	if err = createTables(ctx, db); err != nil {
		return errors.New("can't create tables")
	}
	return nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	const (
		tasksTable = `
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
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

func InsertTask(task pb.Task) (int64, error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return 0, errors.New("can't open db")
	}
	defer db.Close()

	var q = `
	INSERT INTO tasks ( agentid, status, expr, result, begindate) values ($1, $2, $3, $4, $5)
	`
	result, err := db.ExecContext(ctx, q, task.Agentid, task.Status, task.Expr, task.Result, task.Begindate)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateTask(task *pb.Task) error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return errors.New("can't open db")
	}
	defer db.Close()

	var q = `UPDATE tasks SET 
	agentid = $1, 
	status = $2,
	result = $3,
	begindate = $4,
	enddate = $5 				
	WHERE id = $6`
	_, err = db.ExecContext(ctx, q, task.Agentid, task.Status, task.Result, task.Begindate, task.Enddate, task.Id)
	if err != nil {
		fmt.Println("error update: ", err)
		return err
	}

	return nil
}
