package grpc_client

import (
	"context"
	"errors"
	"status-server/logging"
	"status-server/protobuffer_def"

	rgrpc "google.golang.org/grpc"
)

type ComSigClient struct {
	Conn   *rgrpc.ClientConn
	Client protobuffer_def.ComSigServerClient
}

//创建连接池
var comsigClients map[string]*ComSigClient

func init() {
	comsigClients = make(map[string]*ComSigClient)
}

func CreateServiceClient(addr string) (*ComSigClient, error) {
	if addr == "" {
		return nil, errors.New("empty string addr")
	}

	val, ok := comsigClients[addr]
	if ok {
		return val, nil
	}

	conn, err := rgrpc.Dial(addr, rgrpc.WithInsecure())
	if err != nil {
		logging.Logger.Infof("did not connect: %v ", err)
		return nil, err
	}

	sc := protobuffer_def.NewComSigServerClient(conn)
	comSigClient := &ComSigClient{}
	comSigClient.Conn = conn
	comSigClient.Client = sc

	comsigClients[addr] = comSigClient
	logging.Logger.Infof("add %s into client map success", addr)
	return comSigClient, nil
}

func CloseConnections() {

	if comsigClients == nil {
		return
	}

	for _, val := range comsigClients {
		val.Conn.Close()
	}
	comsigClients = nil
}

func (sg *ComSigClient) PostRpcMsg(baseReq *protobuffer_def.BaseRequest) (*protobuffer_def.BaseResponse, error) {
	baseRsp, err := sg.Client.BaseInterface(context.Background(), baseReq)
	if err != nil {
		logging.Logger.Info("query failed, err is ", err)
		return nil, err
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("query failed, code is ", baseRsp.Code)
		return baseRsp, nil
	}

	return baseRsp, nil
}

/*
func Start() {
	r := zookeeper.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:2181"}
		op.Context = context.Background()
		op.Timeout = time.Second * 5
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create GRPC service
	service := grpc.NewService(
		service.Name("test.client2"),
		service.Registry(r),
		service.Context(ctx),
	)

	service.Client().Init(client.Retries(3), client.PoolSize(200), client.PoolTTL(time.Second*20), client.RequestTimeout(time.Second*5))

	test := protobuffer_def.NewStatusServerService("status-service", service.Client())

	for r := 0; r < 20; r++ {
		go func() {
			i := 0
			for {
				_, err := test.BaseInterface(context.Background(), &protobuffer_def.BaseRequest{RequestId: "1111", C: protobuffer_def.CMD_REGISTER_STATUS})
				if err != nil {
					fmt.Println(err)
				} else {
					i++
				}
				if i%10000 == 0 {
					fmt.Println(i, time.Now().Unix())
				}
			}
		}()
	}

}
*/
