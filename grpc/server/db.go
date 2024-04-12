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
	if err = DeleteFromAgents(ctx, db); err != nil {
		return errors.New("error delete from agents")
	}
	return nil
}

func GetTasksFromDb() ([]pb.Task, error) {
	var tasks []pb.Task = make([]pb.Task, 0)
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var q = `SELECT id, agentid, status, expr, result FROM tasks WHERE status = 'in_progress'`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := pb.Task{}
		err := rows.Scan(&task.Id, &task.Agentid, &task.Status, &task.Expr, &task.Result)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
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
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		ip TEXT,
		port INTEGER
	);`

		usersTable = `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		lgn TEXT NOT NULL UNIQUE,
		psw TEXT NOT NULL
	);`
	)

	if _, err := db.ExecContext(ctx, tasksTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, agentsTable); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}
	return nil
}

func GetTaskById(id int64) (Task, error) {
	fmt.Println("from GetTaskById id:", id)
	var task Task = Task{}
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return task, errors.New("can't open db")
	}
	defer db.Close()

	var q = "SELECT id, agentid, status, expr, result, begindate, enddate FROM tasks WHERE id = $1"
	err = db.QueryRowContext(ctx, q, id).Scan(&task.Id, &task.AgentId, &task.Status, &task.Expr, &task.Result, &task.BeginDate, &task.EndDate)
	if err != nil {
		fmt.Println(err)
		return task, err
	}

	fmt.Println("from GetTaskById task:", task)
	return task, err
}

func LoginExists(lgn string) error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return errors.New("can't open db")
	}
	defer db.Close()

	var lgnInTable string
	var q = "SELECT lgn FROM users WHERE lgn = $1"
	err = db.QueryRowContext(ctx, q, lgn).Scan(&lgnInTable)
	if err != nil {
		return err
	}
	if lgn != lgnInTable {
		return errors.New("password invalid")
	}

	return nil
}

func IsPswValid(lgn string, psw string) error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return errors.New("can't open db")
	}
	defer db.Close()

	var pswInTable string
	var q = "SELECT psw FROM users WHERE lgn = $1"
	err = db.QueryRowContext(ctx, q, lgn).Scan(&pswInTable)
	if err != nil {
		return err
	}
	if psw != pswInTable {
		return errors.New("password invalid")
	}

	return nil
}

func InsertUser(lgn string, psw string) error {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return errors.New("can't open db")
	}
	defer db.Close()

	var q = `
	INSERT INTO users ( lgn, psw ) values ($1, $2)
	`
	_, err = db.ExecContext(ctx, q, lgn, psw)
	if err != nil {
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

func DeleteFromAgents(ctx context.Context, db *sql.DB) error {
	var q = `
	DELETE FROM agents
	`
	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

func InsertAgent(agentParams *pb.AgentParams) (int64, error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return 0, errors.New("can't open db")
	}
	defer db.Close()

	var q = `
	INSERT INTO agents (ip, port) values ($1, $2)
	`
	result, err := db.ExecContext(ctx, q, agentParams.Ip, agentParams.Port)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func checkTables() {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		fmt.Println("error open DB")
	}
	defer db.Close()

	var q = `SELECT * FROM tasks`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		fmt.Println("error read rows")
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println("rows")
	}
}
