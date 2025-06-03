package view

type Barrage struct {
	// 剧集id
	VideoSeriesID string `json:"video_series_id"`
	// 浏览器指纹（设备id）
	DeviceID string `json:"device_id"`
	// 弹幕内容
	Content string `json:"content"`
	// 已播放秒数
	PlaySeconds int `json:"play_seconds"`
}

// swagger:parameters GetBarragesByVideoSeriesID
type GetBarragesReq struct {
	// in:header
	Token string `json:"Authorization"`
	// 剧集id
	// in:query
	VsID int `json:"v_s_id"`
}

type BarrageResp struct {
	// id(自增)
	ID int64 `json:"id"`
	// 会员id
	MemberID int64 `json:"member_id"`
	// 浏览器指纹（设备id）
	DeviceID string `json:"device_id"`
	// 弹幕内容
	Content string `json:"content"`
	// 已播放秒数
	PlaySeconds int `json:"play_seconds"`
	// 创建时间
	CreatedAt int64 `json:"created_at"`
}

// swagger:model
type GetBarragesResp struct {
	YsBarrages []*BarrageResp `json:"data"`
	Total      int64          `json:"total"`
}

// swagger:parameters CreateBarrage
type CreateBarrageReqWrap struct {
	// in:header
	Token string `json:"Authorization"`
	// in:body
	Body Barrage
}

// swagger:model
type CreateBarrageResp struct {
	ID int64 `json:"unique_id"`
}

// swagger:parameters DeleteBarrage
type DeleteBarrageReq struct {
	// in:header
	Token string `json:"Authorization"`
	// in:path
	ID int64 `json:"id"`
}
