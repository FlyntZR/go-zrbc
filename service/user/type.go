package service

import (
	"go-zrbc/db"
	"go-zrbc/view"
)

func DBToViewUser(dUser *db.Member) *view.Member {
	vUser := view.Member{
		ID:       dUser.ID,
		User:     dUser.User,
		UserName: dUser.UserName,
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
