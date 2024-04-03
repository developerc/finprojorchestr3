package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/developerc/finprojorchestr3/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

var RegisteredAgentMap map[int]Agent //хранилище зарегистрированных агентов
var IdAgent int                      //счетчик ID агентаvvv
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

/*func (s *Server) SendTask(ctx context.Context, task *pb.Task) (*pb.Task, error) {
	return task, nil
}*/

func (s *Server) PushFinishTask(ctx context.Context, task *pb.Task) (*pb.Task, error) {

	return task, nil
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

// добавить очередь задач и обработчик периодической отсылки задач агенту
// добавить выбор из RegisteredAgentMap очередного агента, отправка задачи//
// если агент не принял, выбирать другого//
func SndTsk(expr string) {
	host := "localhost"
	port := "5001"
	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()

	grpcClient := pb.NewOrchServerServiceClient(conn)
	tskAgent, err := grpcClient.SendTask(context.TODO(), &pb.Task{Id: 1, Expr: expr})
	if err != nil {
		log.Println("failed invoking tskAgent: ", err)
	}
	fmt.Println("tskAgent:  ", tskAgent)
}

func CreateOrchGRPCserver() {
	RegisteredAgentMap = make(map[int]Agent)
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
