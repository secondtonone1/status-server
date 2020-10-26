package constants

const baseKey = "ss" //status-server
const LastTimeKeySubKey = "l"
const NextHeartbeatIntervalSubKey = "n"
const RegisterInfoSubKey = "i"
const userHashKey = "userHashM"
const chatRoomHashKey = "chatRoomHashM"

//取得hmap key ss:{identity}
func GetIdentityKey(identity string) string {
	return baseKey + Colon + identity
}

//最后一次心跳时间 ｛identity｝:l
func GetLastTimeKey(identity string) string {
	return identity + Colon + LastTimeKeySubKey
}

//心跳间隔key  ｛identity｝:n
func GetNextHeartbeatIntervalKey() string {
	return baseKey + Colon + NextHeartbeatIntervalSubKey
}

//注册信息key  ｛identity｝:i
func GetRegisterInfoKey(identity string) string {
	return identity + Colon + RegisterInfoSubKey
}

//获取用户hash表的key
func GetUserHashMKey() string {
	return userHashKey
}

//获取聊天室hash表的key
func GetChatRoomHashMKey() string {
	return chatRoomHashKey
}

//聊天室信息key  ｛chatRoomId｝:i
func GetChatRoomInfoKey(chatRoomId string) string {
	return chatRoomId + Colon + RegisterInfoSubKey
}
