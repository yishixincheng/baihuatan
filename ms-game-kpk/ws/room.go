package ws

import (
	"baihuatan/api/oauth/pb"
	"encoding/json"
	"baihuatan/ms-game-kpk/model"
	"fmt"
	"sync"
	"context"
)

// Room - 房间结构体
type Room struct {

}

// RoomID 房间ID
type RoomID  string

// RoomManager 房间管理者
type RoomManager struct {
	MaxUserCount       int64  // 每个房间最多的用户数
	MaxRoomCount       int64  // 支持最多的房间数
	RoomList           map[RoomID]*Room
	mutex              sync.Mutex 
}

// CreateRoom - 创建房间
func (p *RoomManager) CreateRoom() {

}

// MatchingRoom - 匹配房间
func (p *RoomManager) MatchingRoom() {
	
}

// GetKpkUser - 获取登录者信息
func (p *RoomManager) GetKpkUser(ctx context.Context, userToken string) (*model.KpkUser, error) {
	var userDetail = pb.UserDetails{}
	if err := json.Unmarshal([]byte(userToken), &userDetail); err != nil {
		fmt.Println("err:", err);
		return nil, err
	}

	return model.GetKpkUserByUID(ctx, userDetail.UserID)
}

