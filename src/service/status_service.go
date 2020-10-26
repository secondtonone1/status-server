package service

import (
	"status-server/protobuffer_def"
)

type StatusService interface {
	RegisterStatus(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	QueryStatus(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	CreateChatRoom(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	UpdateOffLine(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	RPCGetChatRoom(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	RPCDelChatRoom(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
	KickPerson(user_id string, curAddr string)
	UpdateOnLine(baseRequest *protobuffer_def.BaseRequest, baseResponse *protobuffer_def.BaseResponse) error
}
