package redis

import (
	"encoding/json"
	"status-server/constants"
	"status-server/grpc_client"
	"status-server/logging"
	"status-server/model"
	"status-server/protobuffer_def"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
)

type statusInfoStruct struct {
	deviceType byte
	statusInfo string
	lastTime   int64
	interval   int32
}

func RegisterStatus(request *protobuffer_def.RegisterStatusRequest) error {

	ud := model.UserData{}
	ud.Avator = request.GetUserAvator()
	ud.Phone = request.GetPhone()
	ud.RegAddr = request.GetRegisterInfo()
	ud.RoomList = request.GetRoomList()
	ud.UserName = request.GetUserName()
	ud.UserId = request.GetIdentity()
	ud.RegAddr = request.GetRegisterInfo()
	hMapKey := constants.GetUserHashMKey()
	userKey := constants.GetIdentityKey(ud.UserId)
	udJson, err := json.Marshal(ud)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return err
	}

	err = Hset(hMapKey, userKey, udJson)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return err
	}
	logging.Logger.Infof("reg status success, hMapKey is %s, userKey is %s, udJson is %s", hMapKey, userKey, udJson)
	return nil
}

func QueryStatus(request *protobuffer_def.QueryStatusRequest) (*protobuffer_def.QueryStatusResponse, error) {

	identify := request.GetIdentity() //唯一标识
	hMapkey := constants.GetUserHashMKey()
	userKey := constants.GetIdentityKey(identify)

	udJson, err := Hget(hMapkey, userKey)
	if err != nil {
		logging.Logger.Info("get user data from redis failed, error is ", err)
		return nil, err
	}

	ud := &model.UserData{}
	err = json.Unmarshal([]byte(udJson), ud)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return nil, err
	}

	rsp := &protobuffer_def.QueryStatusResponse{}
	rsp.Identity = ud.UserId
	rsp.Phone = ud.Phone
	rsp.RoomList = ud.RoomList
	rsp.UserAvator = ud.Avator
	rsp.UserName = ud.UserName
	rsp.State = int32(ud.State)
	rsp.ChatRoomId = ud.ChatRoomId
	rsp.LastChat = ud.LastChat
	rsp.RegAddr = ud.RegAddr
	rsp.OffLine = ud.Offline

	return rsp, nil
}

func CreateChatRoom(request *protobuffer_def.CreateChatRoomReq) (*protobuffer_def.CreateChatRoomRsp, error) {

	curtime := time.Now().UnixNano() / 1e6
	times := strconv.FormatInt(curtime, 10)
	roomId := request.Caller + "-" + request.Answer + "-" + times
	newRoom := &model.ChatRoomData{Caller: request.Caller, Answer: request.Answer, ChatRoomId: roomId}
	chatRoomJs, err := json.Marshal(newRoom)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return nil, err
	}

	//存储chat room信息
	hChatRoomMapkey := constants.GetChatRoomHashMKey()
	chatRoomKey := constants.GetChatRoomInfoKey(roomId)

	hUserMapKey := constants.GetUserHashMKey()
	err = Hset(hChatRoomMapkey, chatRoomKey, chatRoomJs)
	if err != nil {
		logging.Logger.Info("get user data from redis failed, error is ", err)
		return nil, err
	}

	//获取呼叫人key
	userKey := constants.GetIdentityKey(request.Caller)

	udJson, err := Hget(hUserMapKey, userKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", request.Caller, err)
		return nil, err
	}

	ud := &model.UserData{}
	err = json.Unmarshal([]byte(udJson), ud)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return nil, err
	}

	logging.Logger.Info("caller data is ", *ud)
	ud.State = constants.User_Busy
	cur := time.Now().Unix()
	ud.LastChat = cur
	ud.ChatRoomId = roomId

	marshUser, err := json.Marshal(ud)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return nil, err
	}

	err = Hset(hUserMapKey, userKey, marshUser)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return nil, err
	}
	logging.Logger.Infof("create chat room update user info, hMapKey is %s, userKey is %s, udJson is %s", hUserMapKey, userKey, marshUser)

	//获取被叫人key
	userKey = constants.GetIdentityKey(request.Answer)

	udJson, err = Hget(hUserMapKey, userKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", request.Answer, err)
		return nil, err
	}

	ud = &model.UserData{}
	err = json.Unmarshal([]byte(udJson), ud)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return nil, err
	}

	logging.Logger.Info("answer data is ", *ud)
	ud.State = constants.User_Busy
	ud.LastChat = cur
	ud.ChatRoomId = roomId

	marshUser, err = json.Marshal(ud)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return nil, err
	}

	err = Hset(hUserMapKey, userKey, marshUser)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return nil, err
	}

	logging.Logger.Infof("create chat room update user info, hMapKey is %s, userKey is %s, udJson is %s", hUserMapKey, userKey, marshUser)

	rsp := &protobuffer_def.CreateChatRoomRsp{}
	rsp.Answer = request.Answer
	rsp.Caller = request.Caller
	rsp.ChatRoomId = roomId
	return rsp, nil
}

func UpdateOnLine(request *protobuffer_def.UpdateOnline) error {
	hUserMapKey := constants.GetUserHashMKey()
	userKey := constants.GetIdentityKey(request.UserId)

	//更新用户离线状态
	udJson, err := Hget(hUserMapKey, userKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", userKey, err)
		return err
	}

	ud := &model.UserData{}
	err = json.Unmarshal([]byte(udJson), ud)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return err
	}

	logging.Logger.Info("user data is ", *ud)

	if ud.Offline == false {
		return nil
	}

	ud.Offline = false
	ud.RegAddr = request.RegAddr

	marshUser, err := json.Marshal(ud)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return err
	}

	err = Hset(hUserMapKey, userKey, marshUser)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return err
	}

	logging.Logger.Infof("update online success, hMapKey is %s, userKey is %s, udJson is %s", hUserMapKey, userKey, marshUser)

	return nil
}

func UpdateOffLine(request *protobuffer_def.UpdateOfflineReq) error {

	hUserMapKey := constants.GetUserHashMKey()
	userKey := constants.GetIdentityKey(request.Identity)

	//更新用户离线状态
	udJson, err := Hget(hUserMapKey, userKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", request.Identity, err)
		return err
	}

	ud := &model.UserData{}
	err = json.Unmarshal([]byte(udJson), ud)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return err
	}

	logging.Logger.Info("user data is ", *ud)
	roomId := ud.ChatRoomId
	ud.State = constants.User_Idle
	ud.ChatRoomId = ""
	ud.Offline = true
	marshUser, err := json.Marshal(ud)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return err
	}

	err = Hset(hUserMapKey, userKey, marshUser)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return err
	}

	logging.Logger.Infof("update offline success, hMapKey is %s, userKey is %s, udJson is %s", hUserMapKey, userKey, marshUser)

	if roomId == "" {
		logging.Logger.Infof("user %v not in chat room", ud.UserId)
		return nil
	}

	hChatRoomMapkey := constants.GetChatRoomHashMKey()
	chatRoomKey := constants.GetChatRoomInfoKey(roomId)

	//根据聊天房间找到另一个用户
	roomJson, err := Hget(hChatRoomMapkey, chatRoomKey)
	if err != nil {
		logging.Logger.Infof("get chat room data %v from redis failed, error is %v", roomId, err)
		return err
	}

	chatRoomData := &model.ChatRoomData{}
	err = json.Unmarshal([]byte(roomJson), chatRoomData)
	if err != nil {
		logging.Logger.Info("chat room json unmarshal failed, error is ", err)
		return err
	}

	otherUser := chatRoomData.Caller
	if chatRoomData.Caller == request.Identity {
		otherUser = chatRoomData.Answer
	}

	otherKey := constants.GetIdentityKey(otherUser)
	otherJson, err := Hget(hUserMapKey, otherKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", otherUser, err)
		return err
	}

	otherD := &model.UserData{}
	err = json.Unmarshal([]byte(otherJson), otherD)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return err
	}

	logging.Logger.Info("other user data is ", *otherD)

	otherD.State = constants.User_Idle
	otherD.ChatRoomId = ""
	otherD.Offline = false
	marshOther, err := json.Marshal(otherD)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return err
	}

	err = Hset(hUserMapKey, otherKey, marshOther)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return err
	}

	logging.Logger.Infof("update offline success, hMapKey is %s, userKey is %s, otherJson is %s", hUserMapKey, otherKey, marshOther)

	//删除聊天室
	_, err = HDelFiled(hChatRoomMapkey, []string{chatRoomKey})
	if err != nil {
		logging.Logger.Info("del chat room key failed, key is ", hChatRoomMapkey)
	} else {
		logging.Logger.Info("del chat room key success, key is ", hChatRoomMapkey)
	}
	//给另一方服务器发送对端下线通知

	sigClient, err := grpc_client.CreateServiceClient(otherD.RegAddr)
	if err != nil {
		logging.Logger.Info("create sigclient rpc failed, regaddr is ", otherD.RegAddr)
		return nil
	}

	//todo....
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "force_terminate_notify"
	baseReq.C = protobuffer_def.CMD_FORCE_TERMINAL_NOTIFY
	notifyCall := &protobuffer_def.ForceTerminateNotify{}
	notifyCall.ChatRoomId = roomId
	notifyCall.OtherId = otherUser

	body, err := ptypes.MarshalAny(notifyCall)
	if err != nil {
		logging.Logger.Info("proto marshal failed")
		return nil
	}
	baseReq.Body = body

	logging.Logger.Info("single call notify req body is ", notifyCall)

	_, err = sigClient.PostRpcMsg(baseReq)
	if err != nil {
		logging.Logger.Info("post rpc msg to sig failed, otherD addr is ", otherD.RegAddr)
		return nil
	}

	return nil
}

func GetChatRoomData(request *protobuffer_def.GetChatRoomReq) (*protobuffer_def.GetChatRoomRsp, error) {
	hChatRoomMapkey := constants.GetChatRoomHashMKey()
	chatRoomKey := constants.GetChatRoomInfoKey(request.ChatRoomId)

	//根据聊天房间找到另一个用户
	roomJson, err := Hget(hChatRoomMapkey, chatRoomKey)
	if err != nil {
		logging.Logger.Infof("get chat room data %v from redis failed, error is %v", request.ChatRoomId, err)
		return nil, err
	}

	chatRoomData := &model.ChatRoomData{}
	err = json.Unmarshal([]byte(roomJson), chatRoomData)
	if err != nil {
		logging.Logger.Info("chat room json unmarshal failed, error is ", err)
		return nil, err
	}

	chatRsp := &protobuffer_def.GetChatRoomRsp{}
	chatRsp.Answer = chatRoomData.Answer
	chatRsp.Caller = chatRoomData.Caller
	chatRsp.ChatRoomId = chatRoomData.ChatRoomId
	return chatRsp, nil

}

func DelChatRoomData(request *protobuffer_def.DelChatRoomReq) (*protobuffer_def.DelChatRoomRsp, error) {
	hChatRoomMapkey := constants.GetChatRoomHashMKey()
	chatRoomKey := constants.GetChatRoomInfoKey(request.ChatRoomId)

	//根据聊天房间找到另一个用户
	roomJson, err := Hget(hChatRoomMapkey, chatRoomKey)
	if err != nil {
		logging.Logger.Infof("get chat room data %v from redis failed, error is %v", request.ChatRoomId, err)
		return nil, err
	}

	chatRoomData := &model.ChatRoomData{}
	err = json.Unmarshal([]byte(roomJson), chatRoomData)
	if err != nil {
		logging.Logger.Info("chat room json unmarshal failed, error is ", err)
		return nil, err
	}

	//删除聊天室
	_, err = HDelFiled(hChatRoomMapkey, []string{chatRoomKey})
	if err != nil {
		logging.Logger.Info("del chat room key failed, key is ", hChatRoomMapkey)
	} else {
		logging.Logger.Info("del chat room key success, key is ", hChatRoomMapkey)
	}

	hUserMapKey := constants.GetUserHashMKey()
	answerKey := constants.GetIdentityKey(chatRoomData.Answer)
	callerKey := constants.GetIdentityKey(chatRoomData.Caller)

	//更新呼叫方信息
	udJson, err := Hget(hUserMapKey, callerKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", chatRoomData.Caller, err)
		return nil, err
	}

	callerData := &model.UserData{}
	err = json.Unmarshal([]byte(udJson), callerData)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return nil, err
	}

	logging.Logger.Info("caller data is ", *callerData)
	callerData.State = constants.User_Idle
	callerData.ChatRoomId = ""

	marshUser, err := json.Marshal(callerData)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return nil, err
	}

	err = Hset(hUserMapKey, callerKey, marshUser)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return nil, err
	}

	logging.Logger.Infof("update caller success, hMapKey is %s, userKey is %s, udJson is %s", hUserMapKey, callerKey, marshUser)

	//更新被叫方信息
	udJson, err = Hget(hUserMapKey, answerKey)
	if err != nil {
		logging.Logger.Infof("get user data %v from redis failed, error is %v", chatRoomData.Answer, err)
		return nil, err
	}

	answerData := &model.UserData{}
	err = json.Unmarshal([]byte(udJson), answerData)
	if err != nil {
		logging.Logger.Info("json unmarshal failed, error is ", err)
		return nil, err
	}

	logging.Logger.Info("answer data is ", *answerData)
	answerData.State = constants.User_Idle
	answerData.ChatRoomId = ""

	marshUser, err = json.Marshal(answerData)
	if err != nil {
		logging.Logger.Info("json marshal error is : ", err)
		return nil, err
	}

	err = Hset(hUserMapKey, answerKey, marshUser)
	if err != nil {
		logging.Logger.Info("set userdata to redis failed, err is ", err)
		return nil, err
	}

	logging.Logger.Infof("update answer success, hMapKey is %s, userKey is %s, udJson is %s", hUserMapKey, answerKey, marshUser)

	chatRsp := &protobuffer_def.DelChatRoomRsp{}
	return chatRsp, nil

}
