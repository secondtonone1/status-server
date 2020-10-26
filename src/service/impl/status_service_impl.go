package impl

import (
	"fmt"
	"status-server/logging"
	"status-server/protobuffer_def"
	lredis "status-server/redis"
	serivce "status-server/service"
	"sync"

	"status-server/grpc_client"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

var (
	ssi     *statusServiceImpl
	ssiOnce = &sync.Once{}
)

type statusServiceImpl struct{}

func init() {
	fmt.Println("init status")
}

func NewStatusService() serivce.StatusService {
	ssiOnce.Do(func() {
		ssi = &statusServiceImpl{}
	})
	return ssi
}

func (s *statusServiceImpl) preDeal(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse, request proto.Message) bool {
	//解析body
	if baseRequest.GetBody() == nil {
		baseResponse.Code = protobuffer_def.ReturnCode_BODY_IS_NULL
		baseResponse.Desc = "body is null"
		return false
	}
	//反序列化
	err := ptypes.UnmarshalAny(baseRequest.GetBody(), request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "body deserialization error"
		return false
	}
	return true
}

//注册状态
func (s *statusServiceImpl) RegisterStatus(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	//解析请求参数
	request := &protobuffer_def.RegisterStatusRequest{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		logging.Logger.Info("pre deal failed, request is ", baseRequest)
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		return nil
	}

	logging.Logger.Info("receive register status, request is ", request)

	s.KickPerson(request.Identity, request.RegisterInfo)
	//注册状态
	//err := lredis.RegisterStatus(request)
	err := lredis.RegisterStatus(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "save status err"
		return err
	}
	//构建连接池, status->community-sig
	grpc_client.CreateServiceClient(request.RegisterInfo)
	return nil
}

//判断是否执行踢人逻辑
func (s *statusServiceImpl) KickPerson(user_id string, curAddr string) {
	//解析请求参数
	request := &protobuffer_def.QueryStatusRequest{Identity: user_id}

	//缓存中查询用户状态信息
	response, err := lredis.QueryStatus(request)
	if err != nil {
		logging.Logger.Infof("user %v is not in redis ", user_id)
		return
	}

	logging.Logger.Infof("check kick person uid is %v, cur addr is %v,old addr is %v", user_id, curAddr, response.RegAddr)
	oldAddr := response.RegAddr

	if oldAddr == "" {
		logging.Logger.Infof("oldAddr is empty")
		return
	}

	if curAddr == oldAddr {
		logging.Logger.Infof("cur addr %v is same as %v", curAddr, response.RegAddr)
		return
	}
	//离线跳过踢人
	if response.OffLine == true {
		logging.Logger.Infof("cur addr %v is same as %v", curAddr, response.RegAddr)
		return
	}

	logging.Logger.Infof("cur addr %v is not same as %v", curAddr, response.RegAddr)

	//todo....
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "kick_user"
	baseReq.C = protobuffer_def.CMD_CMD_KICK_USER
	notifyCall := &protobuffer_def.KickPerson{}
	notifyCall.UserId = user_id

	body, err := ptypes.MarshalAny(notifyCall)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return
	}
	baseReq.Body = body

	logging.Logger.Info("kick person notify req body is ", notifyCall)

	sigClient, err := grpc_client.CreateServiceClient(oldAddr)
	if err != nil {
		logging.Logger.Info("create sigclient rpc failed, regaddr is ", oldAddr)
		return
	}

	_, err = sigClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("post rpc msg to sig failed, old addr is ", oldAddr)
		return
	}

}

func (s *statusServiceImpl) UpdateOnLine(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	//解析请求参数
	request := &protobuffer_def.UpdateOnline{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		logging.Logger.Info("pre deal failed, request is ", baseRequest)
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		return nil
	}

	logging.Logger.Info("update online  request is ", request)

	s.KickPerson(request.UserId, request.RegAddr)
	//注册状态
	//err := lredis.RegisterStatus(request)
	err := lredis.UpdateOnLine(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "save status err"
		return err
	}
	return nil
}

//查询状态
func (s *statusServiceImpl) QueryStatus(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS

	//解析请求参数
	request := &protobuffer_def.QueryStatusRequest{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		return nil
	}

	//缓存中查询用户状态信息
	response, err := lredis.QueryStatus(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "save status error"
		return err
	}

	//序列化状态信息
	if response != nil {
		body, err2 := ptypes.MarshalAny(response)
		if err2 != nil {
			baseResponse.Code = protobuffer_def.ReturnCode_SERIALIZATION_ERROR
			baseResponse.Desc = "MarshalAny error"
			return err2
		}
		baseResponse.Body = body
		return err2
	}
	return nil
}

//创建聊天室
func (s *statusServiceImpl) CreateChatRoom(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS

	//解析请求参数
	request := &protobuffer_def.CreateChatRoomReq{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		return nil
	}

	//缓存中查询用户状态信息
	response, err := lredis.CreateChatRoom(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "save status error"
		return err
	}

	//序列化状态信息
	if response != nil {
		body, err2 := ptypes.MarshalAny(response)
		if err2 != nil {
			baseResponse.Code = protobuffer_def.ReturnCode_SERIALIZATION_ERROR
			baseResponse.Desc = "MarshalAny error"
			return err2
		}
		baseResponse.Body = body
		return err2
	}
	return nil
}

//更新用户离线信息
func (s *statusServiceImpl) UpdateOffLine(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error {
	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS

	//解析请求参数
	request := &protobuffer_def.UpdateOfflineReq{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		return nil
	}

	//缓存中查询用户状态信息
	err := lredis.UpdateOffLine(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "update offline error"
		return err
	}

	return nil
}

//rpc 服务收到获取聊天室请求
func (s *statusServiceImpl) RPCGetChatRoom(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.GetChatRoomReq{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "json unmarshal failed"
		return nil
	}

	//缓存中查询用户状态信息
	roomrsp, err := lredis.GetChatRoomData(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "get chat room error"
		return err
	}

	//序列化状态信息
	if roomrsp != nil {
		body, err2 := ptypes.MarshalAny(roomrsp)
		if err2 != nil {
			baseResponse.Code = protobuffer_def.ReturnCode_SERIALIZATION_ERROR
			baseResponse.Desc = "MarshalAny error"
			return err2
		}
		baseResponse.Body = body
		return err2
	}
	return nil
}

func (s *statusServiceImpl) RPCDelChatRoom(baseRequest *protobuffer_def.BaseRequest,
	baseResponse *protobuffer_def.BaseResponse) error {

	baseResponse.Code = protobuffer_def.ReturnCode_SUCCESS
	baseResponse.Desc = "success"
	//解析请求参数
	request := &protobuffer_def.DelChatRoomReq{}
	if !s.preDeal(baseRequest, baseResponse, request) {
		baseResponse.Code = protobuffer_def.ReturnCode_DESERIALIZATION_ERROR
		baseResponse.Desc = "json unmarshal failed"
		return nil
	}

	//删除缓存中房间信息
	roomrsp, err := lredis.DelChatRoomData(request)
	if err != nil {
		baseResponse.Code = protobuffer_def.ReturnCode_UNKOWN_ERROR
		baseResponse.Desc = "del chat room error"
		return err
	}

	//序列化状态信息
	if roomrsp != nil {
		body, err2 := ptypes.MarshalAny(roomrsp)
		if err2 != nil {
			baseResponse.Code = protobuffer_def.ReturnCode_SERIALIZATION_ERROR
			baseResponse.Desc = "MarshalAny error"
			return err2
		}
		baseResponse.Body = body
		return err2
	}
	return nil
}
