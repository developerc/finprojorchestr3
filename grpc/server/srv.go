package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/developerc/finprojorchestr3/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

var RegisteredAgentMap map[int]Agent  //хранилище зарегистрированных агентов
var RegisteredTaskMap map[int]pb.Task //хранилище задач
var TaskQueue []pb.Task               //очередь задач
var IdAgent int                       //счетчик ID агента
var IdTask int                        //счетчик Id задач
var mutex sync.Mutex

type Agent struct {
	Id   int    `json:"id"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	//mutex sync.Mutex
}

type Server struct {
	pb.OrchServerServiceServer // сервис из сгенерированного пакета
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) HBreq(ctx context.Context, heartBit *pb.HeartBit) (*pb.HeartBitResp, error) {
	//hbr := pb.HeartBitResp
	return &pb.HeartBitResp{}, nil
}

// обрабатываем запрос на регистрацию агента
func (s *Server) RegisterNewAgent(ctx context.Context, in *pb.AgentParams) (*pb.AgentParamsResponse, error) {
	var agent Agent = Agent{}

	mutex.Lock()
	IdAgent++
	agent.Id = IdAgent
	agent.Ip = in.Ip
	agent.Port = int(in.Port)
	RegisteredAgentMap[IdAgent] = agent
	mutex.Unlock()
	log.Println("RegisteredAgentMap: ", RegisteredAgentMap)
	return &pb.AgentParamsResponse{Id: int32(IdAgent)}, nil
}

func (s *Server) PushFinishTask(ctx context.Context, task *pb.Task) (*pb.Task, error) {
	fmt.Println("принимаем решенную задачу: ", task)
	//TODO сделать занесение решенной задачи в БД
	UpdateTask(task)
	return task, nil
}

// добавить очередь задач и обработчик периодической отсылки задач агенту
func HandleHttpExpr(expr string) {
	var task pb.Task
	mutex.Lock()
	//IdTask++
	//task.Id = int32(IdTask)
	task.Expr = expr
	task.Status = "start"
	task.Begindate = time.Now().Unix()

	id, err := InsertTask(task)
	if err != nil {
		log.Println("could not insert task: ", err)
	}
	task.Id = int32(id)
	RegisteredTaskMap[int(task.Id)] = task
	TaskQueue = append(TaskQueue, task)
	mutex.Unlock()
	//fmt.Println(TaskQueue)
}

// добавить выбор из RegisteredAgentMap очередного агента, отправка задачи//
// если агент не принял, выбирать другого//
// обработчик очереди задач
func handlerTaskQueue() {
	for {
		if len(TaskQueue) > 0 { //если в очереди есть задачи, начинаем работу
			if len(RegisteredAgentMap) > 0 { //если есть зарегистрированные агенты
				//fmt.Println(TaskQueue, RegisteredAgentMap)
				for _, agent := range RegisteredAgentMap {
					if tskAgent, err := SndTsk(agent, &TaskQueue[0]); tskAgent != nil {
						if err != nil {
							log.Println("could not send task: ", err)
							continue
						}
						mutex.Lock()
						RegisteredTaskMap[int(tskAgent.Id)] = *tskAgent
						mutex.Unlock()
						TaskQueue = TaskQueue[1:]
						UpdateTask(tskAgent)
						break
					}
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func SndTsk(agent Agent, task *pb.Task) (*pb.Task, error) {
	host := agent.Ip                         //"localhost"
	port := strconv.Itoa(agent.Port)         //"5001"
	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		//os.Exit(1)
		return nil, err
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()

	grpcClient := pb.NewOrchServerServiceClient(conn)
	tskAgent, err := grpcClient.SendTask(context.TODO(), task /*&pb.Task{Id: 1, Expr: expr}*/)
	if err != nil {
		log.Println("failed invoking tskAgent: ", err)
		return nil, err
	}
	fmt.Println("tskAgent:  ", tskAgent)

	return tskAgent, nil
}

func CreateOrchGRPCserver() {
	CreateSqliteDb()
	RegisteredAgentMap = make(map[int]Agent)
	RegisteredTaskMap = make(map[int]pb.Task)
	go handlerTaskQueue()
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию серверной части OrchServerServiceServer
	orchserverServiceServer := NewServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterOrchServerServiceServer(grpcServer, orchserverServiceServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
