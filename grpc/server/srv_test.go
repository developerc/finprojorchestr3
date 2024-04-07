package server

import (
	"context"
	"database/sql"

	//"fmt"
	"testing"

	pb "github.com/developerc/finprojorchestr3/proto"
)

func TestCreateSqliteDb(t *testing.T) {
	CreateSqliteDb()
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		t.Error("error open DB")
	}
	defer db.Close()

	var q = `SELECT * FROM tasks`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		t.Error("error read tasks")
	}

	q = `SELECT * FROM agents`
	rows, err = db.QueryContext(ctx, q)
	if err != nil {
		t.Error("error read agents")
	}
	defer rows.Close()
}

func TestInsertTask(t *testing.T) {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		t.Error("error open DB")
	}
	defer db.Close()

	var q = `
	DELETE FROM tasks
	`
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		t.Error("error TestInsertTask, can't delete from tasks")
	}

	var task pb.Task
	_, err = InsertTask(task)
	if err != nil {
		t.Error("error TestInsertTask, can't insert into tasks")
	}
}

func TestInsertAgent(t *testing.T) {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		t.Error("error open DB")
	}
	defer db.Close()

	var agentParams pb.AgentParams
	agentParams.Ip = "test"
	id, err := InsertAgent(&agentParams)
	if err != nil {
		t.Error("error TestInsertTask, can't insert into agents")
	}
	//fmt.Println("id: ", id)
	var q = "SELECT id, ip, port FROM agents WHERE id = $1"
	err = db.QueryRowContext(ctx, q, id).Scan(&agentParams.Id, &agentParams.Ip, &agentParams.Port)
	if err != nil {
		t.Error("error select from agents")
	}
	//fmt.Println(agentParams)
	if agentParams.Id != int32(id) || agentParams.Ip != "test" {
		t.Error("error insert into agents")
	}
}

// тест проводится при запущенном агенте
func TestRegisterNewAgent(t *testing.T) {
	RegisteredAgentMap = make(map[int]Agent)
	ctx := context.TODO()
	var agentParams pb.AgentParams
	agentParams.Ip = "test"
	agentParams.Port = 6000
	_, err := NewServer().RegisterNewAgent(ctx, &agentParams)
	if err != nil {
		t.Error("error register new agent")
	}
}

func TestHandleHttpExpr(t *testing.T) {
	RegisteredTaskMap = make(map[int]pb.Task)
	HandleHttpExpr("3+2")
	if len(RegisteredTaskMap) == 0 {
		t.Error("error handle http expression")
	}
}

/*func TestUpdateTask(t *testing.T) {
	var task pb.Task
	task.Expr = "test"
	id, err := InsertTask(task)
	if err != nil {
		t.Error("error TestInsertTask, can't insert into tasks")
	}
	fmt.Println("id: ", id)
	task.Id = int32(id)
	task.Expr = "proba"
	err = UpdateTask(&task)
	if err != nil {
		t.Error("error update tasks")
	}

}*/
