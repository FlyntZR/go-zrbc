package wschannel

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"go-zrbc/pkg/xlog"
	"go-zrbc/view"
)

var maxUserInRoom int64 = 1000

type Room struct {
	MaxUserInRoom int64
	ID            string // 房间id
	RType         int    // 房间类型：0 影视；1 直播
	// Total users in the room currently
	Total int64

	sync.RWMutex
	// userID: client
	clients map[string]*Client
}

func NewRoom(roomID string, rType int) *Room {
	r := &Room{
		MaxUserInRoom: maxUserInRoom,
		ID:            roomID,
		RType:         rType,
		clients:       make(map[string]*Client, 0),
	}
	return r
}

func (r *Room) Close() error {
	atomic.StoreInt64(&r.Total, 0)
	r.RLock()
	defer r.RUnlock()
	for _, cli := range r.clients {
		cli.Close("close")
	}
	return nil
}

func (r *Room) TotalClients() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.clients)
}

func (r *Room) GetClientByConnID(connID string) *Client {
	r.RLock()
	defer r.RUnlock()

	return r.clients[connID]
}

func (r *Room) BroadcastToAllClients(cli *Client, wsmsg *view.WsResp) error {
	r.RLock()
	defer r.RUnlock()
	msg, _ := json.Marshal(wsmsg)
	for _, c := range r.clients {
		xlog.Infof("BroadcastToAllClients cli:%+v", c)
		select {
		case c.bytesSend <- msg:
		default:
		}
	}
	return nil
}

func (r *Room) Desc() int64 {
	return int64(atomic.AddInt64(&r.Total, -1))
}

func (r *Room) Incr() int64 {
	return int64(atomic.AddInt64(&r.Total, 1))
}

func (r *Room) AddClient(cli *Client) error {
	// local room user number limit
	total := r.Incr()
	if total > r.MaxUserInRoom {
		return ErrRoomFull
	}
	r.Lock()
	defer r.Unlock()
	r.clients[cli.connID] = cli

	return nil
}

func (r *Room) RemoveClient(cli *Client) {
	r.Lock()
	defer r.Unlock()

	delete(r.clients, cli.ConnID())
	if len(r.clients) == 0 {
		delete(cli.mgr.rooms, cli.Room.ID)
	}
	r.Desc()
}
