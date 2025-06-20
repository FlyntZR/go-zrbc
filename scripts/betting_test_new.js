import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// Configuration - 移到最前面
const WS_IP = __ENV.WS_IP || 'ws://192.168.0.213';
const WS_URL = `${WS_IP}/15109`;
const WS_URL_15101 = `${WS_IP}/15101`;
const ACCOUNT_COUNT = Number(__ENV.ACCOUNT_COUNT) || 2;
const ACCOUNTS = Array.from({ length: ACCOUNT_COUNT }, (_, i) => `laugh_g_${i + 1}`);
const DEBUG = __ENV.DEBUG === 'true'; // Debug flag to control logging

// Helper function for conditional logging
function debugLog(message) {
  if (DEBUG) {
    console.log(message);
  }
}

function debugError(message) {
  if (DEBUG) {
    console.error(message);
  }
}

export const options = {
  vus: Number(__ENV.ACCOUNT_COUNT) || 2,
  duration: '5m',
};

export default function () {
  const account = ACCOUNTS[__VU % ACCOUNTS.length];
  
  // Connect to 15109 for authentication
  const authRes = ws.connect(WS_URL, {}, function (socket) {
    socket.on('open', function () {
      debugLog(`连接到 ${WS_URL} 账号: ${account}`);
      
      // Send authentication message
      const authMsg = JSON.stringify({
        protocol: 0,
        data: {
          account: account,
          password: "123456",
          dtBetLimitSelectID: {},
          bGroupList: false,
          videoName: "TC",
          videoDelay: 3000,
          userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
        }
      });
      
      socket.send(authMsg);
    });

    socket.on('message', function (message) {
      try {
        const data = JSON.parse(message);
        
        if (data.protocol === 0) {
          // Authentication successful, get SID
          const sid = data.data.sid;
          debugLog(`15109登录成功, sid: ${sid}`);
          
          // Connect to 15101 with the SID
          connectTo15101(sid, account);
          
          // Close the authentication socket
          socket.close();
        }
      } catch (e) {
        debugError(`15109解析消息失败: ${e.message}, 消息内容: ${message}`);
      }
    });

    socket.on('error', function (e) {
      console.error(`15109连接错误 ${account}:`, e);
    });

    socket.on('close', function () {
      console.log(`15109连接关闭 ${account}`);
    });
  });

  const authCheck = check(authRes, { '15109连接成功': (r) => r && r.status === 101 });
  if (!authCheck['15109连接成功']) {
    console.error(`15109连接失败 ${account}, 状态: ${authRes.status}`);
  }
}

function connectTo15101(sid, account) {
  console.log(`开始连接15101, 账号: ${account}, sid: ${sid}`);
  
  const res = ws.connect(WS_URL_15101, {}, function (socket) {
    socket.on('open', function () {
      console.log(`15101 WebSocket连接已打开: ${WS_URL_15101} 账号: ${account}`);
      // Send login message with SID
      const loginMsg = JSON.stringify({
        protocol: 1,
        data: {
          dtBetLimitSelectID: {
            "101": 124, "102": 125, "103": 9, "104": 126, "105": 127,
            "106": 128, "107": 129, "108": 149, "110": 131, "111": 150,
            "112": 250, "113": 251, "117": 260, "121": 261, "125": 600,
            "126": 599, "128": 584, "129": 602, "301": 29
          },
          bGroupList: false,
          videoName: "TC",
          videoDelay: 3000,
          userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
          sid: sid
        }
      });
      socket.send(loginMsg);
      debugLog(`发送15101登录消息: ${loginMsg}`);
    });

    let firstTwentyOneFlag = true;
    let modifyBetLimitSendFlag = false;
    let modifyBetLimitFlag = false;
    let recvBetResultFlag = false;
    let betSerialNumber = 1;
    let groupID = 0;
    let memberID = 0;

    socket.on('message', function (message) {
      try {
        const data = JSON.parse(message);
        if (data.protocol === 0) {
          console.log(`15101 登录成功: ${JSON.stringify(message)}, account: ${account}`);
        } else if (data.protocol === 10) {
          // Join table response
          if (data.data && data.data.memberID) {
            memberID = data.data.memberID;
          }
          console.log(`15101 进入桌台成功: ${JSON.stringify(message)}`);
          // Send modify bet limit request
          const modifyLimitMsg = JSON.stringify({
            protocol: 60,
            data: {
              dtBetLimitSelectID: {
                "101": 2, "102": 125, "103": 9, "104": 126, "105": 127,
                "106": 128, "107": 129, "108": 149, "110": 131, "111": 150,
                "112": 250, "113": 251, "117": 260, "121": 261, "125": 600,
                "126": 599, "128": 584, "129": 602, "301": 29
              }
            }
          });
          socket.send(modifyLimitMsg);
          modifyBetLimitSendFlag = true;
          debugLog(`15101 修改限红`);
        } else if (data.protocol === 60) {
          // Modify bet limit response
          if (data.data && data.data.memberID && memberID !== 0 && data.data.memberID !== memberID) return;
          if (modifyBetLimitSendFlag) {
              modifyBetLimitFlag = true;
              debugLog(`15101 修改限红成功: ${JSON.stringify(message)}`);
          }
        } else if (data.protocol === 25) {
          // Game result
          if (data.data && data.data.groupID !== groupID && groupID !== 0) return;
          debugLog(`15101 得到一局结果: ${JSON.stringify(message)}, account: ${account}`);
          if (modifyBetLimitFlag) {
            recvBetResultFlag = true;
          }
        } else if (data.protocol === 31) {
          // Payout
          if (data.data && data.data.memberID && memberID !== 0 && data.data.memberID !== memberID) {
              debugLog(`15101 派彩失败: ${JSON.stringify(message)}`);
              debugLog(`15101 派彩失败 groupID: ${groupID}, memberID: ${memberID}`);
              return;
          }
          console.log(`15101 派彩成功: account: ${account} ${JSON.stringify(message)}`);
        } else if (data.protocol === 38) {
          // Bet time
          const betTimeData = data.data;
          if (betTimeData.groupID !== groupID && groupID !== 0) return;
          if (firstTwentyOneFlag) {
            groupID = betTimeData.groupID;
            firstTwentyOneFlag = false;
            const joinTableMsg = JSON.stringify({
              protocol: 10,
              data: {
                dtBetLimitSelectID: {
                  "101": 124, "102": 125, "103": 9, "104": 126, "105": 127,
                  "106": 128, "107": 129, "108": 149, "110": 131, "111": 150,
                  "112": 250, "113": 251, "117": 260, "121": 261, "125": 600,
                  "126": 599, "128": 584, "129": 602, "301": 29
                },
                groupID: groupID
              }
            });
            socket.send(joinTableMsg);
            debugLog(`15101 发送登录桌台: ${joinTableMsg}, groupID: ${groupID}`);
            debugLog(`15101 得到下注时间剩余秒数信息第一次: ${JSON.stringify(betTimeData)}`);
            return;
          }
          debugLog(`15101 得到下注时间剩余秒数: ${JSON.stringify(betTimeData)}`);
          if (recvBetResultFlag && modifyBetLimitFlag) {
            // Send betting message
            const bettingMsg = JSON.stringify({
              protocol: 22,
              data: {
                betSerialNumber: betSerialNumber,
                gameNo: betTimeData.gameNo,
                gameNoRound: betTimeData.gameNoRound,
                betArr: [{
                  betArea: [1, 2, 3][randomIntBetween(0, 2)],
                  addBetMoney: [100, 200, 300][randomIntBetween(0, 2)]
                }],
                commission: 0
              }
            });
            socket.send(bettingMsg);
            debugLog('15101 投注成功: ' + bettingMsg + ', account: ' + account);
            recvBetResultFlag = false;
            betSerialNumber++;
          }
        } else if (data.protocol === 22) {
          // Betting response
          if (data.data && data.data.groupID !== groupID && groupID !== 0) {
              debugLog(`15101 投注成功回复失败: ${JSON.stringify(message)}`);
              debugLog(`15101 投注成功回复失败 groupID: ${groupID}, memberID: ${memberID}`);
              return;
          }
          debugLog(`15101 投注成功回复: ${JSON.stringify(message)}`);
        } else {
        }
      } catch (e) {
        debugError(`15101解析消息失败: ${e.message}, 消息内容: ${message}`);
      }
    });

    socket.on('error', function (e) {
      console.error(`15101连接错误 ${account}:`, e);
    });

    socket.on('close', function () {
      console.log(`15101连接关闭 ${account}`);
    });

    // Send ping every 5 seconds
    const pingInterval = setInterval(function () {
      socket.ping();
    }, 5000);

    // Clean up interval on close
    socket.on('close', function () {
      clearInterval(pingInterval);
    });
  });

  const checkResult = check(res, { '15101连接成功': (r) => r && r.status === 101 });
  if (!checkResult['15101连接成功']) {
    console.error(`15101连接失败 ${account}, 状态: ${res.status}`);
  }
} 