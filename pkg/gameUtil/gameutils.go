package gameUtil

import (
	"strconv"
	"strings"
)

// ChipsCheck contains valid chip values
var ChipsCheck = []int64{
	1, 5, 10, 20, 50, 100, 500, 1000, 5000, 10000, 20000, 50000, 100000, 200000, 1000000,
	5000000, 10000000, 20000000, 50000000, 10000000, 20000000, 50000000, 100000000,
}

// ReportFormat formats a row of game data into a map
func ReportFormat(row map[string]interface{}, lang string) map[string]interface{} {
	tmp := make(map[string]interface{})
	tmp["user"] = row["username"]
	tmp["betId"] = row["bet01"]
	tmp["betTime"] = row["bet08"]
	tmp["beforeCash"] = row["bet12"]
	tmp["bet"] = row["bet13"]
	tmp["validbet"] = row["bet41"]
	tmp["water"] = row["bet16"]
	tmp["result"] = row["bet17"]
	tmp["betResult"] = GetBetContent(row["bet02"].(string), row["bet09"].(string), lang)
	tmp["waterbet"] = row["bet41"]
	tmp["winLoss"] = row["bet14"].(float64) - row["bet13"].(float64)
	tmp["ip"] = row["ip"]
	tmp["gid"] = row["bet02"]
	tmp["event"] = row["bet03"]
	tmp["eventChild"] = row["bet04"]
	tmp["round"] = row["bet03"]
	tmp["subround"] = row["bet04"]
	tmp["tableId"] = row["bet39"]
	tmp["commission"] = row["commission"]
	tmp["settime"] = row["updatetime"]
	tmp["reset"] = row["bet38"]
	tmp["gameResult"] = row["gameResult"]
	tmp["gname"] = GetLangText(row["cnname"].(string), lang)
	return tmp
}

// GetBetContent converts bet content based on game type
func GetBetContent(gtype, content, lang string) string {
	switch gtype {
	case "101", "102", "301", "126":
		betName := map[string]string{
			"Dragon":            GetLangText("龙", lang),
			"DragonOdd":         GetLangText("龙单", lang),
			"DragonEven":        GetLangText("龙双", lang),
			"DragonRed":         GetLangText("龙红", lang),
			"DragonBlack":       GetLangText("龙黑", lang),
			"Tiger":             GetLangText("虎", lang),
			"TigerOdd":          GetLangText("虎单", lang),
			"TigerEven":         GetLangText("虎双", lang),
			"TigerRed":          GetLangText("虎红", lang),
			"TigerBlack":        GetLangText("虎黑", lang),
			"Tie":               GetLangText("和", lang),
			"Banker":            GetLangText("庄", lang),
			"Player":            GetLangText("闲", lang),
			"BPair":             GetLangText("庄对", lang),
			"PPair":             GetLangText("闲对", lang),
			"Small":             GetLangText("小", lang),
			"Big":               GetLangText("大", lang),
			"BankerNatural":     GetLangText("庄例牌", lang),
			"PlayerNatural":     GetLangText("闲例牌", lang),
			"Super6":            GetLangText("幸运6", lang),
			"AnyPair":           GetLangText("任意对子", lang),
			"PerfectPair":       GetLangText("完美对子", lang),
			"BankerDragonBonus": GetLangText("庄龙宝", lang),
			"PlayerDragonBonus": GetLangText("闲龙宝", lang),
			"SuperTie0":         GetLangText("超和0点", lang),
			"SuperTie1":         GetLangText("超和1点", lang),
			"SuperTie2":         GetLangText("超和2点", lang),
			"SuperTie3":         GetLangText("超和3点", lang),
			"SuperTie4":         GetLangText("超和4点", lang),
			"SuperTie5":         GetLangText("超和5点", lang),
			"SuperTie6":         GetLangText("超和6点", lang),
			"SuperTie7":         GetLangText("超和7点", lang),
			"SuperTie8":         GetLangText("超和8点", lang),
			"SuperTie9":         GetLangText("超和9点", lang),
			"Tip_1_":            GetLangText("小费", lang),
		}
		return betName[content]

	case "103":
		switch content {
		case "Odd":
			return GetLangText("单", lang)
		case "Even":
			return GetLangText("双", lang)
		case "Big":
			return GetLangText("大", lang)
		case "Small":
			return GetLangText("小", lang)
		case "Red":
			return GetLangText("红", lang)
		case "Black":
			return GetLangText("黑", lang)
		default:
			if strings.HasPrefix(content, "Num") {
				parts := strings.Split(content, "Num")
				nums := strings.Split(parts[1], "_")
				if len(nums) == 1 {
					return GetLangText("单号", lang) + ":" + nums[0]
				}
				return GetLangText("多号", lang) + ":" + strings.Join(nums, ",")
			} else if strings.HasPrefix(content, "Dozen") {
				parts := strings.Split(content, "Dozen")
				nums := strings.Split(parts[1], "_")
				return GetLangText("区块", lang) + ":" + nums[0] + GetLangText("至", lang) + nums[1]
			} else if strings.HasPrefix(content, "Column") {
				parts := strings.Split(content, "Column")
				return GetLangText("直排", lang) + ":" + parts[1]
			} else if content == "Tip_1_" {
				return GetLangText("小费", lang)
			}
			return content
		}

	case "104":
		switch content {
		case "Odd":
			return GetLangText("单", lang)
		case "Even":
			return GetLangText("双", lang)
		case "Big":
			return GetLangText("大", lang)
		case "Small":
			return GetLangText("小", lang)
		case "AllLeopard":
			return GetLangText("全围", lang)
		default:
			if strings.HasPrefix(content, "Leopard") {
				parts := strings.Split(content, "Leopard")
				return parts[1] + GetLangText("围", lang)
			} else if strings.HasPrefix(content, "TwoDice") {
				parts := strings.Split(content, "TwoDice")
				n1 := parts[1][0:1]
				n2 := parts[1][1:2]
				return GetLangText("双骰", lang) + ":" + n1 + "," + n2
			} else if strings.HasPrefix(content, "Sum") {
				parts := strings.Split(content, "Sum")
				return GetLangText("点数合", lang) + ":" + parts[1]
			} else if strings.HasPrefix(content, "OneDice") {
				parts := strings.Split(content, "OneDice")
				return GetLangText("单骰", lang) + ":" + parts[1]
			} else if strings.HasPrefix(content, "DoubleDice") {
				parts := strings.Split(content, "DoubleDice")
				return GetLangText("对子", lang) + ":" + parts[1]
			} else if content == "Tip_1_" {
				return GetLangText("小费", lang)
			}
			return content
		}

	case "105":
		betName := map[string]string{
			"Player1Equal":  GetLangText("闲一平倍", lang),
			"Player2Equal":  GetLangText("闲二平倍", lang),
			"Player3Equal":  GetLangText("闲三平倍", lang),
			"Player1Double": GetLangText("闲一翻倍", lang),
			"Player2Double": GetLangText("闲二翻倍", lang),
			"Player3Double": GetLangText("闲三翻倍", lang),
			"Tip_1_":        GetLangText("小费", lang),
		}
		return betName[content]

	case "106", "112", "113":
		betName := map[string]string{
			"BankerPairPlus":   GetLangText("庄对牌以上", lang),
			"Player1Win":       GetLangText("闲1赢", lang),
			"Player1Lose":      GetLangText("闲1输", lang),
			"Player1Tie":       GetLangText("闲1和", lang),
			"Player1ThreeFace": GetLangText("闲1三公", lang),
			"Player1Pair":      GetLangText("闲1对子", lang),
			"Player1PairPlus":  GetLangText("闲1对牌以上", lang),
			"Player2Win":       GetLangText("闲2赢", lang),
			"Player2Lose":      GetLangText("闲2输", lang),
			"Player2Tie":       GetLangText("闲2和", lang),
			"Player2ThreeFace": GetLangText("闲2三公", lang),
			"Player2Pair":      GetLangText("闲2对子", lang),
			"Player2PairPlus":  GetLangText("闲2对牌以上", lang),
			"Player3Win":       GetLangText("闲3赢", lang),
			"Player3Lose":      GetLangText("闲3输", lang),
			"Player3Tie":       GetLangText("闲3和", lang),
			"Player3ThreeFace": GetLangText("闲3三公", lang),
			"Player3Pair":      GetLangText("闲3对子", lang),
			"Player3PairPlus":  GetLangText("闲3对牌以上", lang),
			"Tip_1_":           GetLangText("小费", lang),
		}
		return betName[content]

	case "107":
		betName := map[string]string{
			"Odd":    GetLangText("单", lang),
			"Even":   GetLangText("双", lang),
			"Fan1":   GetLangText("1番", lang),
			"Fan2":   GetLangText("2番", lang),
			"Fan3":   GetLangText("3番", lang),
			"Fan4":   GetLangText("4番", lang),
			"Nim12":  GetLangText("1念2", lang),
			"Nim13":  GetLangText("1念3", lang),
			"Nim14":  GetLangText("1念4", lang),
			"Nim21":  GetLangText("2念1", lang),
			"Nim23":  GetLangText("2念3", lang),
			"Nim24":  GetLangText("2念4", lang),
			"Nim31":  GetLangText("3念1", lang),
			"Nim32":  GetLangText("3念2", lang),
			"Nim34":  GetLangText("3念4", lang),
			"Nim43":  GetLangText("4念3", lang),
			"Nim42":  GetLangText("4念2", lang),
			"Nim41":  GetLangText("4念1", lang),
			"Kwok12": GetLangText("12角", lang),
			"Kwok23": GetLangText("23角", lang),
			"Kwok34": GetLangText("34角", lang),
			"Kwok41": GetLangText("14角", lang),
			"Kwok14": GetLangText("14角", lang),
			"Nga231": GetLangText("23一通", lang),
			"Nga241": GetLangText("24一通", lang),
			"Nga341": GetLangText("34一通", lang),
			"Nga132": GetLangText("13二通", lang),
			"Nga142": GetLangText("14二通", lang),
			"Nga342": GetLangText("34二通", lang),
			"Nga143": GetLangText("14三通", lang),
			"Nga243": GetLangText("24三通", lang),
			"Nga123": GetLangText("12三通", lang),
			"Nga124": GetLangText("12四通", lang),
			"Nga134": GetLangText("13四通", lang),
			"Nga234": GetLangText("23四通", lang),
			"Ssh234": GetLangText("三门234", lang),
			"Ssh134": GetLangText("三门134", lang),
			"Ssh124": GetLangText("三门124", lang),
			"Ssh123": GetLangText("三门123", lang),
			"Tip_1_": GetLangText("小费", lang),
		}
		return betName[content]

	case "108":
		betName := map[string]string{
			"Big":    GetLangText("大", lang),
			"Small":  GetLangText("小", lang),
			"Odd":    GetLangText("单", lang),
			"Even":   GetLangText("双", lang),
			"R4":     GetLangText("四红", lang),
			"W4":     GetLangText("四白", lang),
			"R3W1":   GetLangText("三红一白", lang),
			"W3R1":   GetLangText("三白一红", lang),
			"Tip_1_": GetLangText("小费", lang),
		}
		return betName[content]

	case "110":
		betName := map[string]string{
			"Odd":        GetLangText("单", lang),
			"Even":       GetLangText("双", lang),
			"Big":        GetLangText("大", lang),
			"Small":      GetLangText("小", lang),
			"Triples1":   GetLangText("围骰：鱼", lang),
			"Triples2":   GetLangText("围骰：虾", lang),
			"Triples3":   GetLangText("围骰：葫芦", lang),
			"Triples4":   GetLangText("围骰：钱币", lang),
			"Triples5":   GetLangText("围骰：蟹", lang),
			"Triples6":   GetLangText("围骰：鸡", lang),
			"Anytriples": GetLangText("全围", lang),
			"Anycolor":   GetLangText("任意三色", lang),
			"Color1r":    GetLangText("指定单色：红", lang),
			"Color1g":    GetLangText("指定单色：绿", lang),
			"Color1b":    GetLangText("指定单色：蓝", lang),
			"Color2r":    GetLangText("指定双色：红", lang),
			"Color2g":    GetLangText("指定双色：绿", lang),
			"Color2b":    GetLangText("指定双色：蓝", lang),
			"Color3r":    GetLangText("指定三色：红", lang),
			"Color3g":    GetLangText("指定三色：绿", lang),
			"Color3b":    GetLangText("指定三色：蓝", lang),
			"Dice1":      GetLangText("三军：鱼", lang),
			"Dice2":      GetLangText("三军：虾", lang),
			"Dice3":      GetLangText("三军：葫芦", lang),
			"Dice4":      GetLangText("三军：钱币", lang),
			"Dice5":      GetLangText("三军：蟹", lang),
			"Dice6":      GetLangText("三军：鸡", lang),
			"Sum4":       GetLangText("点数总和4", lang),
			"Sum5":       GetLangText("点数总和5", lang),
			"Sum6":       GetLangText("点数总和6", lang),
			"Sum7":       GetLangText("点数总和7", lang),
			"Sum8":       GetLangText("点数总和8", lang),
			"Sum9":       GetLangText("点数总和9", lang),
			"Sum10":      GetLangText("点数总和10", lang),
			"Sum11":      GetLangText("点数总和11", lang),
			"Sum12":      GetLangText("点数总和12", lang),
			"Sum13":      GetLangText("点数总和13", lang),
			"Sum14":      GetLangText("点数总和14", lang),
			"Sum15":      GetLangText("点数总和15", lang),
			"Sum16":      GetLangText("点数总和16", lang),
			"Sum17":      GetLangText("点数总和17", lang),
			"Tip_1_":     GetLangText("小费", lang),
		}
		return betName[content]

	case "111":
		betName := map[string]string{
			"Dragon":        GetLangText("龙", lang),
			"Phoenix":       GetLangText("凤", lang),
			"Pair9Plus":     GetLangText("对九以上", lang),
			"Straight":      GetLangText("顺子", lang),
			"Flush":         GetLangText("同花", lang),
			"StraightFlush": GetLangText("同花顺", lang),
			"AnyTriple":     GetLangText("豹子", lang),
			"Tip_1_":        GetLangText("小费", lang),
		}
		return betName[content]

	case "117", "121":
		betName := map[string]string{
			"Ante":     GetLangText("底注", lang),
			"Bet":      GetLangText("加注", lang),
			"PairPlus": GetLangText("奖赏", lang),
		}
		return betName[content]

	default:
		return content
	}
}

// GetGameResultString converts game result to readable string
func GetGameResultString(gtype, result string) string {
	// Handle cancel and shuffle cases
	if result == "cancel" {
		return GetLangText("该局取消", "cn")
	}
	if result == "0" {
		return GetLangText("洗牌...", "cn")
	}

	// Handle different game types
	switch gtype {
	case "101", "301": // Baccarat
		if strings.HasPrefix(result, ";") {
			result = result[1:]
		}
		resultary := strings.Split(result, ";")
		if len(resultary) < 2 {
			return result
		}
		resultary1 := strings.Split(resultary[0], ":")
		resultary2 := strings.Split(resultary[1], ":")
		if len(resultary1) < 2 || len(resultary2) < 2 {
			return result
		}
		resultary3 := strings.Split(resultary1[1], ",")
		resultary4 := strings.Split(resultary2[1], ",")
		cards := append(resultary3, resultary4...)
		for i, card := range cards {
			if card != "" {
				cards[i] = GetPokerCardFlower(parseInt(card))
			}
		}
		return GetLangText("庄", "cn") + ":" + cards[0] + cards[1] + cards[2] + " " + GetLangText("闲", "cn") + ":" + cards[3] + cards[4] + cards[5]

	case "102", "126": // Dragon Tiger
		if strings.HasPrefix(result, ";") {
			result = result[1:]
		}
		resultary := strings.Split(result, ";")
		if len(resultary) < 2 {
			return result
		}
		resultary1 := strings.Split(resultary[0], ":")
		resultary2 := strings.Split(resultary[1], ":")
		if len(resultary1) < 2 || len(resultary2) < 2 {
			return result
		}
		resultary3 := strings.Split(resultary1[1], ",")
		resultary4 := strings.Split(resultary2[1], ",")
		cards := append(resultary3, resultary4...)
		for i, card := range cards {
			if card != "" {
				cards[i] = GetPokerCardFlower(parseInt(card))
			}
		}
		return GetLangText("龙", "cn") + ":" + cards[0] + " " + GetLangText("虎", "cn") + ":" + cards[1]

	case "103": // Roulette
		return result

	case "104": // Sicbo
		return result

	case "105": // Three Card Poker
		resultary := strings.Split(result, ";")
		if len(resultary) < 5 {
			return result
		}
		aryBanker := strings.Split(resultary[1], ":")
		aryPlayer1 := strings.Split(resultary[2], ":")
		aryPlayer2 := strings.Split(resultary[3], ":")
		aryPlayer3 := strings.Split(resultary[4], ":")

		niucard := []string{
			getValueOrDefault(aryBanker, 1, ""),
			getValueOrDefault(aryPlayer1, 1, ""),
			getValueOrDefault(aryPlayer2, 1, ""),
			getValueOrDefault(aryPlayer3, 1, ""),
		}

		var card [][]string
		for _, value := range niucard {
			reCard := strings.Split(value, ",")
			var cardRow []string
			for _, r := range reCard {
				if r != "" {
					cardRow = append(cardRow, GetPokerCardFlower(parseInt(r)))
				}
			}
			card = append(card, cardRow)
		}

		return GetLangText("庄", "cn") + ":" + strings.Join(card[0], ",") + " " +
			GetLangText("闲1", "cn") + ":" + strings.Join(card[1], ",") + " " +
			GetLangText("闲2", "cn") + ":" + strings.Join(card[2], ",") + " " +
			GetLangText("闲3", "cn") + ":" + strings.Join(card[3], ",")

	case "106": // Bull Bull
		resultary := strings.Split(result, ";")
		if len(resultary) < 5 {
			return result
		}
		aryBanker := strings.Split(resultary[1], ":")
		aryPlayer1 := strings.Split(resultary[2], ":")
		aryPlayer2 := strings.Split(resultary[3], ":")
		aryPlayer3 := strings.Split(resultary[4], ":")

		niucard := []string{
			getValueOrDefault(aryBanker, 1, ""),
			getValueOrDefault(aryPlayer1, 1, ""),
			getValueOrDefault(aryPlayer2, 1, ""),
			getValueOrDefault(aryPlayer3, 1, ""),
		}

		var card [][]string
		for _, value := range niucard {
			reCard := strings.Split(value, ",")
			var cardRow []string
			for _, r := range reCard {
				if r != "" {
					cardRow = append(cardRow, GetPokerCardFlower(parseInt(r)))
				}
			}
			card = append(card, cardRow)
		}

		return GetLangText("庄", "cn") + ":" + card[0][0] + card[0][1] + card[0][2] + " " +
			GetLangText("闲1", "cn") + ":" + card[1][0] + card[1][1] + card[1][2] + " " +
			GetLangText("闲2", "cn") + ":" + card[2][0] + card[2][1] + card[2][2] + " " +
			GetLangText("闲3", "cn") + ":" + card[3][0] + card[3][1] + card[3][2]

	case "107": // Mahjong
		return result

	case "108": // Texas Hold'em
		switch result {
		case "0":
			return GetLangText("四白", "cn")
		case "1":
			return GetLangText("三白一红", "cn")
		case "2":
			return GetLangText("二紅二白", "cn")
		case "3":
			return GetLangText("三红一白", "cn")
		case "4":
			return GetLangText("四红", "cn")
		default:
			return result
		}

	case "110": // Fish
		return result

	case "111": // Three Card Brag
		if strings.HasPrefix(result, ";") {
			result = result[1:]
		}
		resultary := strings.Split(result, ";")
		if len(resultary) < 2 {
			return result
		}
		resultary1 := strings.Split(resultary[0], ":")
		resultary2 := strings.Split(resultary[1], ":")
		if len(resultary1) < 2 || len(resultary2) < 2 {
			return result
		}
		resultary3 := strings.Split(resultary1[1], ",")
		resultary4 := strings.Split(resultary2[1], ",")
		cards := append(resultary3, resultary4...)
		for i, card := range cards {
			if card != "" {
				cards[i] = GetPokerCardFlower(parseInt(card))
			}
		}
		return GetLangText("龙", "cn") + ":" + cards[0] + cards[1] + cards[2] + "  " + GetLangText("凤", "cn") + ":" + cards[3] + cards[4] + cards[5]

	case "112": // Blackjack
		resultary := strings.Split(result, ";")
		if len(resultary) < 5 {
			return result
		}
		aryBanker := strings.Split(resultary[1], ":")
		aryPlayer1 := strings.Split(resultary[2], ":")
		aryPlayer2 := strings.Split(resultary[3], ":")
		aryPlayer3 := strings.Split(resultary[4], ":")

		niucard := []string{
			getValueOrDefault(aryBanker, 1, ""),
			getValueOrDefault(aryPlayer1, 1, ""),
			getValueOrDefault(aryPlayer2, 1, ""),
			getValueOrDefault(aryPlayer3, 1, ""),
		}

		var card [][]string
		for _, value := range niucard {
			reCard := strings.Split(value, ",")
			var cardRow []string
			for _, r := range reCard {
				if r != "" {
					cardRow = append(cardRow, GetPokerCardFlower(parseInt(r)))
				}
			}
			card = append(card, cardRow)
		}

		return GetLangText("庄", "cn") + ":" + card[0][0] + card[0][1] + "," + " " +
			GetLangText("闲1", "cn") + ":" + card[1][0] + card[1][1] + "," + " " +
			GetLangText("闲2", "cn") + ":" + card[2][0] + card[2][1] + "," + " " +
			GetLangText("闲3", "cn") + ":" + card[3][0] + card[3][1]

	case "113": // Baccarat Squeeze
		cards := map[string]string{
			"101": GetLangText("1筒", "cn"),
			"102": GetLangText("2筒", "cn"),
			"103": GetLangText("3筒", "cn"),
			"104": GetLangText("4筒", "cn"),
			"105": GetLangText("5筒", "cn"),
			"106": GetLangText("6筒", "cn"),
			"107": GetLangText("7筒", "cn"),
			"108": GetLangText("8筒", "cn"),
			"109": GetLangText("9筒", "cn"),
			"137": GetLangText("白板", "cn"),
		}

		var tmp []string
		resultary := strings.Split(result, ";")
		if len(resultary) > 1 {
			for _, res := range resultary[1:] {
				carset := strings.Split(res, ":")
				if len(carset) < 2 {
					continue
				}
				str := ""
				switch carset[0] {
				case "b":
					str += GetLangText("庄", "cn")
				case "p1":
					str += GetLangText("闲1", "cn")
				case "p2":
					str += GetLangText("闲2", "cn")
				case "p3":
					str += GetLangText("闲3", "cn")
				}
				str += ":"
				flowers := strings.Split(carset[1], ",")
				for _, v := range flowers {
					if card, ok := cards[v]; ok {
						str += card
					}
				}
				tmp = append(tmp, str)
			}
		}
		return strings.Join(tmp, ",")

	case "117": // Three Card Poker Squeeze
		if strings.HasPrefix(result, ";") {
			result = result[1:]
		}
		resultary := strings.Split(result, ";")
		if len(resultary) < 2 {
			return result
		}
		resultary1 := strings.Split(resultary[0], ":")
		resultary2 := strings.Split(resultary[1], ":")
		if len(resultary1) < 2 || len(resultary2) < 2 {
			return result
		}
		resultary3 := strings.Split(resultary1[1], ",")
		resultary4 := strings.Split(resultary2[1], ",")
		cards := append(resultary3, resultary4...)
		for i, card := range cards {
			if card != "" {
				cards[i] = GetPokerCardFlower(parseInt(card))
			}
		}
		return GetLangText("庄", "cn") + ":" + strings.Join(cards[0:5], ",") + " " +
			GetLangText("闲", "cn") + ":" + strings.Join(cards[5:10], ",")

	case "121": // Baccarat Squeeze
		if strings.HasPrefix(result, ";") {
			result = result[1:]
		}
		resultary := strings.Split(result, ";")
		if len(resultary) < 2 {
			return result
		}
		resultary1 := strings.Split(resultary[0], ":")
		resultary2 := strings.Split(resultary[1], ":")
		if len(resultary1) < 2 || len(resultary2) < 2 {
			return result
		}
		resultary3 := strings.Split(resultary1[1], ",")
		resultary4 := strings.Split(resultary2[1], ",")
		cards := append(resultary3, resultary4...)
		for i, card := range cards {
			if card != "" {
				cards[i] = GetPokerCardFlower(parseInt(card))
			}
		}
		return GetLangText("庄", "cn") + ":" + strings.Join(cards[0:3], ",") + " " +
			GetLangText("闲", "cn") + ":" + strings.Join(cards[3:6], ",")
	}

	return result
}

// Helper function to parse string to int
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Helper function to get value from slice with default
func getValueOrDefault(slice []string, index int, defaultValue string) string {
	if index < len(slice) {
		return slice[index]
	}
	return defaultValue
}

// GetPokerCardFlower converts poker card number to readable format
func GetPokerCardFlower(card int) string {
	cardNumber := card % 20
	switch cardNumber {
	case 1:
		cardNumber = 0 // A
	case 11:
		cardNumber = 0 // J
	case 12:
		cardNumber = 0 // Q
	case 13:
		cardNumber = 0 // K
	}

	cardFlower := card / 20
	var flower string
	switch cardFlower {
	case 0:
		flower = "♣"
	case 1:
		flower = "♦"
	case 2:
		flower = "♥"
	case 3:
		flower = "♠"
	default:
		flower = ""
	}

	var number string
	switch cardNumber {
	case 0:
		number = "A"
	case 11:
		number = "J"
	case 12:
		number = "Q"
	case 13:
		number = "K"
	default:
		number = strconv.Itoa(cardNumber)
	}

	return flower + number
}
