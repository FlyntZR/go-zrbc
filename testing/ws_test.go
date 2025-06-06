package main

import (
	"encoding/json"
	"go-zrbc/view"
	"log"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestYMZR_ws_15109(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	stopCh := make(chan struct{})

	// c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8082/wss?vsid=5193479026126848", nil)
	// c, _, err := websocket.DefaultDialer.Dial("wss://8822-37-157-223-27.ngrok-free.app/15109", nil)
	c, _, err := websocket.DefaultDialer.Dial("wss://a45gs-t.wmetg.com/15109", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Printf("connect success!")

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"protocol":0,"data":{"account":"laugh","password":"123456","dtBetLimitSelectID":{},"bGroupList":false,"videoName":"TC","videoDelay":3000,"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}}`))
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5 * time.Second) // send a unsolicited pong frame every 15 seconds
	begin_time := time.Now()
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			close(stopCh)
			ticker.Stop()
			return
		case <-ticker.C:
			log.Printf("%v sending a unsolicited pong frame, has been running for %v seconds", time.Now().Format(time.RFC3339), time.Since(begin_time).Seconds())
			err = c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Fatal(err)
			}
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Error reading due to ", err)
			} else {
				txt := string(message)
				log.Printf("recv: %s", txt)
			}
		}
	}
}

type AuthResp struct {
	// bOk
	BOk bool `json:"bOk"`
	// sid
	Sid string `json:"sid"`
}

func connect_ws_15101(ch chan string) {
	sid := <-ch

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	stopCh := make(chan struct{})

	c, _, err := websocket.DefaultDialer.Dial("wss://a45gs-t.wmetg.com/15101", nil)
	if err != nil {
		log.Fatal("dial 15101:", err)
	}
	defer c.Close()
	log.Printf("connect 15101 success!")

	jsonData := []byte(`{"protocol":1,"data":{"dtBetLimitSelectID":{"101":124,"102":125,"103":9,"104":126,"105":127,"106":128,"107":129,"108":149,"110":131,"111":150,"112":250,"113":251,"117":260,"121":261,"125":600,"126":599,"128":584,"129":602,"301":29},"bGroupList":false,"videoName":"TC","videoDelay":3000,"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}}`)

	// Step 1: Parse the JSON into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Step 2: Add the new field to the "data" object
	if data, ok := jsonMap["data"].(map[string]interface{}); ok {
		data["sid"] = sid
	} else {
		log.Fatalf("Failed to access 'data' field")
	}

	// Step 3: Convert back to []byte
	updatedJSON, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = c.WriteMessage(websocket.TextMessage, updatedJSON)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5 * time.Second) // send a unsolicited pong frame every 15 seconds
	begin_time := time.Now()
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			close(stopCh)
			ticker.Stop()
			return
		case <-ticker.C:
			log.Printf("%v sending a unsolicited pong frame, has been running for %v seconds", time.Now().Format(time.RFC3339), time.Since(begin_time).Seconds())
			err = c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Fatal(err)
			}
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Error reading due to ", err)
			} else {
				var wsData view.WsResp
				err = json.Unmarshal(message, &wsData)
				if err != nil {
					log.Fatal("Error Unmarshal to ", err)
				}
				if wsData.Protocol == 1 {
					log.Printf("15101 登录成功: %s", message)
					ticker.Stop()
					return
				}
				log.Printf("15101 recv: %s", message)
			}
		}
	}
}

func TestYMZR_ws_all(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	stopCh := make(chan struct{})

	c, _, err := websocket.DefaultDialer.Dial("wss://a45gs-t.wmetg.com/15109", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Printf("connect 15109 success!")

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"protocol":0,"data":{"account":"laugh","password":"123456","dtBetLimitSelectID":{},"bGroupList":false,"videoName":"TC","videoDelay":3000,"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}}`))
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个字符串类型的通道
	ch := make(chan string)
	go connect_ws_15101(ch)

	ticker := time.NewTicker(5 * time.Second) // send a unsolicited pong frame every 15 seconds
	begin_time := time.Now()
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			close(stopCh)
			ticker.Stop()
			return
		case <-ticker.C:
			log.Printf("%v sending a unsolicited pong frame, has been running for %v seconds", time.Now().Format(time.RFC3339), time.Since(begin_time).Seconds())
			err = c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Fatal(err)
			}
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Error reading due to ", err)
			} else {
				// txt := string(message)
				// log.Printf("recv: %s", txt)
				// bb, err := json.Marshal(message)
				// if err != nil {
				// 	log.Fatal("Error Marshal to ", err)
				// }
				var wsData view.WsResp
				err = json.Unmarshal(message, &wsData)
				if err != nil {
					log.Fatal("Error Unmarshal to ", err)
				}
				// log.Printf("recv: %v", wsData)
				if wsData.Protocol == 0 {
					cc, err := json.Marshal(wsData.Data)
					if err != nil {
						log.Fatal("Error Marshal to ", err)
					}
					var authResp AuthResp
					err = json.Unmarshal(cc, &authResp)
					if err != nil {
						log.Fatal("Error Unmarshal to ", err)
					}
					log.Printf("15109登录成功, sid: %s", authResp.Sid)
					ch <- authResp.Sid
					ticker.Stop()
					break
				}
			}
		}
	}
}

func exec_ws_betting_15101(chbetInfo <-chan view.WsBettingCh) {
	for msg := range chbetInfo {
		time.Sleep(1 * time.Second)

		wsReq, err := json.Marshal(view.WsReq{
			Protocol: 22,
			Data:     msg.BetCh,
		})
		if err != nil {
			log.Fatal("Error Marshal to ", err)
		}
		err = msg.Conn.WriteMessage(websocket.TextMessage, wsReq)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("15101 bet success: %s", wsReq)
	}
}

func connect_ws_betting_15101(ch chan string) {
	sid := <-ch

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	betCh := make(chan view.WsBettingCh)
	defer func() {
		close(betCh)
	}()
	go exec_ws_betting_15101(betCh)

	c, _, err := websocket.DefaultDialer.Dial("wss://59f7-37-157-223-27.ngrok-free.app/15101", nil)
	if err != nil {
		log.Fatal("dial 15101:", err)
	}
	defer c.Close()
	log.Printf("connect 15101 success!")

	jsonData := []byte(`{"protocol":1,"data":{"dtBetLimitSelectID":{"101":124,"102":125,"103":9,"104":126,"105":127,"106":128,"107":129,"108":149,"110":131,"111":150,"112":250,"113":251,"117":260,"121":261,"125":600,"126":599,"128":584,"129":602,"301":29},"bGroupList":false,"videoName":"TC","videoDelay":3000,"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}}`)

	// Step 1: Parse the JSON into a map
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Step 2: Add the new field to the "data" object
	if data, ok := jsonMap["data"].(map[string]interface{}); ok {
		data["sid"] = sid
	} else {
		log.Fatalf("Failed to access 'data' field")
	}

	// Step 3: Convert back to []byte
	updatedJSON, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = c.WriteMessage(websocket.TextMessage, updatedJSON)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5 * time.Second) // send a unsolicited pong frame every 5 seconds
	// periodicMsgTicker := time.NewTicker(2 * time.Second) // send a betting frame every 25 seconds
	begin_time := time.Now()
	defer func() {
		ticker.Stop()
	}()

	recvBetFlag := true
	firstTwentyOneFlag := true
	joinTableFlag := false
	betSerialNumber := 1
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			return
		case <-ticker.C:
			log.Printf("%v sending a unsolicited ping frame, has been running for %v seconds", time.Now().Format(time.RFC3339), time.Since(begin_time).Seconds())
			err = c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Fatal(err)
			}
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Error reading due to ", err)
			} else {
				var wsData view.WsResp
				err = json.Unmarshal(message, &wsData)
				if err != nil {
					log.Fatal("Error Unmarshal to ", err)
				}
				if wsData.Protocol == 0 {
					log.Printf("15101 登录成功: %s", message)
				} else if wsData.Protocol == 22 {
					var WsBettingResp view.WsBettingResp
					err = json.Unmarshal(message, &WsBettingResp)
					if err != nil {
						log.Fatal("Error Unmarshal to ", err)
					}
					if WsBettingResp.Data.GroupID != 3 {
						continue
					} else {
						log.Printf("15101 投注成功回复: %s", message)
						recvBetFlag = true
					}
				} else if wsData.Protocol == 25 {
					log.Printf("15101 得到一局结果: %s", message)
				} else if wsData.Protocol == 31 {
					log.Printf("15101 派彩成功: %s", message)
				} else if wsData.Protocol == 10 {
					log.Printf("15101 进入桌台成功: %s", message)
					joinTableFlag = true
				} else if wsData.Protocol == 38 {
					var wsBetTimeResp view.WsBetTimeResp
					err = json.Unmarshal(message, &wsBetTimeResp)
					if err != nil {
						log.Fatal("Error Unmarshal to ", err)
					}
					if wsBetTimeResp.Data.GroupID != 3 {
						continue
					}
					log.Printf("15101 得到下注时间剩余秒数: %s", message)
					if firstTwentyOneFlag {
						firstTwentyOneFlag = false
						err = c.WriteMessage(websocket.TextMessage, []byte(`{"protocol":10,"data":{"dtBetLimitSelectID":{"101":124,"102":125,"103":9,"104":126,"105":127,"106":128,"107":129,"108":149,"110":131,"111":150,"112":250,"113":251,"117":260,"121":261,"125":600,"126":599,"128":584,"129":602,"301":29},"groupID":3}}`))
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("15101 join table success")
						log.Printf("15101 得到下注时间剩余秒数信息第一次: %v", wsBetTimeResp)
						continue
					}
					log.Printf("15101 得到下注时间剩余秒数: %v", wsBetTimeResp)
					if recvBetFlag && joinTableFlag {
						var wsBettingData = view.WsBettingData{
							BetSerialNumber: betSerialNumber,
							GameNo:          wsBetTimeResp.Data.GameNo,
							GameNoRound:     wsBetTimeResp.Data.GameNoRound,
							BetArr: []view.WsBettingInfoItem{
								{
									BetArea:     1,
									AddBetMoney: 100,
								},
							},
							Commission: 0,
						}
						betCh <- view.WsBettingCh{
							Conn:  c,
							BetCh: wsBettingData,
						}
						recvBetFlag = false
						betSerialNumber++
					}
				} else {
					log.Printf("15101 recv: %s", "other msg")
				}
			}
		}
	}
}

func TestYMZR_ws_betting_all(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial("wss://59f7-37-157-223-27.ngrok-free.app/15109", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Printf("connect 15109 success!")

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"protocol":0,"data":{"account":"laugh","password":"123456","dtBetLimitSelectID":{},"bGroupList":false,"videoName":"TC","videoDelay":3000,"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}}`))
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个字符串类型的通道
	ch := make(chan string)
	go connect_ws_betting_15101(ch)

	ticker := time.NewTicker(5 * time.Second) // send a unsolicited pong frame every 15 seconds
	begin_time := time.Now()
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			ticker.Stop()
			return
		case <-ticker.C:
			log.Printf("%v sending a unsolicited ping frame, has been running for %v seconds", time.Now().Format(time.RFC3339), time.Since(begin_time).Seconds())
			err = c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Fatal(err)
			}
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatal("Error reading due to ", err)
			} else {
				// txt := string(message)
				// log.Printf("recv: %s", txt)
				// bb, err := json.Marshal(message)
				// if err != nil {
				// 	log.Fatal("Error Marshal to ", err)
				// }
				var wsData view.WsResp
				err = json.Unmarshal(message, &wsData)
				if err != nil {
					log.Fatal("Error Unmarshal to ", err)
				}
				// log.Printf("recv: %v", wsData)
				if wsData.Protocol == 0 {
					cc, err := json.Marshal(wsData.Data)
					if err != nil {
						log.Fatal("Error Marshal to ", err)
					}
					var authResp AuthResp
					err = json.Unmarshal(cc, &authResp)
					if err != nil {
						log.Fatal("Error Unmarshal to ", err)
					}
					log.Printf("15109登录成功, sid: %s", authResp.Sid)
					ch <- authResp.Sid
					ticker.Stop()
					break
				}
			}
		}
	}
}
