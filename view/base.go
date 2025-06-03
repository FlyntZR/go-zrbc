package view

// swagger:parameters GetTimeTs
type GetTimeTsReq struct {
	// in:header
	Token string `json:"Authorization"`
}

// swagger:model
type GetTimeTsResp struct {
	STm string `json:"s_tm"`
	Tm  string `json:"tm"`
	Ts  int64  `json:"ts"`
	Tz  string `json:"tz"`
}
