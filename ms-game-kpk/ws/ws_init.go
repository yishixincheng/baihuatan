package ws

// RoomM 创建房间管理器
var RoomM *RoomManager

// InitWs 初始化变量
func InitWs() {
	RoomM = &RoomManager{
		MaxRoomCount : 10000,
		RoomList: make(map[RoomIDType]*Room),
	}
}