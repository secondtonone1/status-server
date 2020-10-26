package model

type UserData struct {
	UserId     string   `json:"userId"`     //用户id
	UserName   string   `json:"userName"`   //用户名
	Phone      string   `json:"phone"`      //电话
	Avator     string   `json:"avator"`     //头像
	RoomList   []string `json:"roomList"`   //房间列表
	State      int      `json:"state"`      //状态
	ChatRoomId string   `json:"chatRoomId"` //聊天室id
	LastChat   int64    `json:"lastChat"`   //上次会话时间
	RegAddr    string   `json:"regAddr"`    //服务器注册地址--"ip:端口"
	Offline    bool     `json:"offLine"`    //是否离线
}

type ChatRoomData struct {
	ChatRoomId string `json:"chatRoomId"` //聊天室id
	Caller     string `json:"caller"`     //呼叫方caller
	Answer     string `json:"answer"`     //被叫方answer
}
