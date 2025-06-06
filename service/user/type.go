package service

import (
	"go-zrbc/db"
	"go-zrbc/view"
)

func DBToViewUser(dUser *db.Member) *view.Member {
	vUser := view.Member{
		ID:              dUser.ID,
		User:            dUser.User,
		UserName:        dUser.UserName,
		Mem003:          dUser.Password,
		Mem005:          dUser.Mem005.Unix(),
		Mem006:          dUser.Mem006,
		Mem007:          dUser.Mem007,
		Mem008:          dUser.Mem008,
		Mem009:          dUser.Mem009,
		Mem010:          dUser.Mem010,
		Mem011:          dUser.Mem011,
		Mem012:          dUser.Mem012,
		Mem013:          dUser.Mem013.Unix(),
		Mem014:          dUser.Mem014,
		Mem015:          dUser.Mem015,
		Mem016:          dUser.Mem016,
		Mem017:          dUser.Mem017,
		Mem018:          dUser.Mem018,
		Mem019:          dUser.Mem019,
		Mem020:          dUser.Mem020,
		Mem021:          dUser.Mem021,
		Mem022:          dUser.Mem022,
		Mem022a:         dUser.Mem022a,
		Mem023:          dUser.Mem023,
		Mem024:          dUser.Mem024,
		Mem026:          dUser.Mem026,
		Mem028:          dUser.Mem028,
		Type:            dUser.Type,
		Currency:        dUser.Currency,
		Cash:            dUser.Cash,
		Money:           dUser.Money,
		Lockmoney:       dUser.Lockmoney,
		Head:            dUser.Head,
		Chips:           dUser.Chips,
		Follow1:         dUser.Follow1,
		Follow2:         dUser.Follow2,
		Tip:             dUser.Tip,
		Red:             dUser.Red,
		Wallet:          dUser.Wallet,
		Opengame:        dUser.Opengame,
		Site:            dUser.Site,
		Lineid:          dUser.Lineid,
		Kickperiod:      dUser.Kickperiod,
		Identity:        dUser.Identity,
		Singlebetprompt: dUser.Singlebetprompt,
		Conwinprompt:    dUser.Conwinprompt,
		Winlossprompt:   dUser.Winlossprompt,
		Onlineprompt:    dUser.Onlineprompt,
		Profitprompt:    dUser.Profitprompt,
	}
	return &vUser
}

func DBToViewUserCache(dUser *db.Member) *view.MemberCache {
	vUser := view.MemberCache{
		UID:      dUser.ID,
		Account:  dUser.User,
		Password: dUser.Password,
		Enable:   dUser.Mem016,
		Mem007:   dUser.Mem007,
		Mem008:   dUser.Mem008,
		Mem009:   dUser.Mem009,
		Mem010:   dUser.Mem010,
		Mem011:   dUser.Mem011,
		SN:       dUser.Mem011,
		Name:     dUser.UserName,
		ULV:      dUser.Mem006,
		ENS:      dUser.Mem016,
		LogFail:  dUser.Mem015,
		UTP:      "M",
	}
	return &vUser
}

func DBToViewAgent(dAgent *db.Agent) *view.Agent {
	vAgent := view.Agent{
		ID:           dAgent.Age001,
		VendorID:     dAgent.Age002,
		Name:         dAgent.Age003,
		Password:     dAgent.Age004,
		CreatedAt:    dAgent.Age005.Unix(),
		RiskReset:    dAgent.Age014,
		Status:       dAgent.Age015,
		BetStatus:    dAgent.Age016,
		ReportSwitch: dAgent.Age019,
		ReportFormat: dAgent.Age020,
		ReportLang:   dAgent.Age021,
		TestLine:     dAgent.Age022,
		Remark:       dAgent.Age024,
		Credit:       dAgent.Credit.String(),
		Cash:         dAgent.Cash.String(),
		Type:         dAgent.Type,
		Currency:     dAgent.Currency,
		Tip:          dAgent.Tip,
		Red:          dAgent.Red,
		PrefixAdd:    dAgent.PrefixAdd,
		PrefixAcc:    dAgent.PrefixAcc,
		ReceiptAcc:   dAgent.ReceiptAcc,
		Notification: dAgent.Notification,
		Sacc:         dAgent.Sacc,
		Opengame:     dAgent.Opengame,
		Site:         dAgent.Site,
		Membermax:    dAgent.Membermax,
		Profitmax:    dAgent.Profitmax,
		Promotecode:  dAgent.Promotecode,
		Kickperiod:   dAgent.Kickperiod,
		Identity:     dAgent.Identity,
		ChkKey:       dAgent.ChkKey,
		ChkLock:      dAgent.ChkLock,
		LastChpsw:    dAgent.LastChpsw.Unix(),
	}
	return &vAgent
}

func DBToViewAgentsLoginPass(dAgentsLoginPass *db.AgentsLoginPass) *view.AgentsLoginPass {
	vAgentsLoginPass := view.AgentsLoginPass{
		ID:           dAgentsLoginPass.ID,
		Aid:          dAgentsLoginPass.Aid,
		VendorID:     dAgentsLoginPass.VendorID,
		Signature:    dAgentsLoginPass.Signature,
		Signature2:   dAgentsLoginPass.Signature2,
		Password:     dAgentsLoginPass.Password,
		Addtime:      dAgentsLoginPass.Addtime.Unix(),
		URL:          dAgentsLoginPass.URL,
		Skyname:      dAgentsLoginPass.Skyname,
		Type:         dAgentsLoginPass.Type,
		Lang:         dAgentsLoginPass.Lang,
		Betfeedback:  dAgentsLoginPass.Betfeedback,
		GatewayURL:   dAgentsLoginPass.GatewayURL,
		WhiteList:    dAgentsLoginPass.WhiteList,
		Operator:     dAgentsLoginPass.Operator,
		ModifyTime:   dAgentsLoginPass.ModifyTime.Unix(),
		Object:       dAgentsLoginPass.Object,
		Settle:       dAgentsLoginPass.Settle,
		Co:           dAgentsLoginPass.Co,
		PrefixSwitch: dAgentsLoginPass.PrefixSwitch,
		OpenGameURL:  dAgentsLoginPass.OpenGameURL,
		Subdomain:    dAgentsLoginPass.Subdomain,
	}
	return &vAgentsLoginPass
}
