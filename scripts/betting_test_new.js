import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

/*
使用示例:
k6 run -e WS_IP=ws://192.168.0.213 -e ACCOUNT_COUNT=5 -e GROUP_IDS=1001,1002,1003,1004,1005 -e DEBUG=true -e OUTPUT_FILE=test_results.json -e DURATION=10m -e PAYOUT_COUNT=10 scripts/betting_test_new.js

参数说明:
- WS_IP: WebSocket服务器地址
- ACCOUNT_COUNT: 账号数量
- GROUP_IDS: 可用的groupID列表，用逗号分隔（如：1001,1002,1003）
- DEBUG: 是否开启调试日志（true/false）
- OUTPUT_FILE: 测试结果输出文件路径（可选，如：test_results.json）
- DURATION: 测试最大执行时间（可选，如：10m, 1h, 30s）
- PAYOUT_COUNT: 每个账号需要完成的派彩次数（可选，默认10次）

每个账号会从GROUP_IDS数组中随机选择一个groupID作为自己的目标桌台。
测试会在以下任一条件满足时结束：
1. 所有账号都完成指定次数的派彩（由PAYOUT_COUNT参数控制）
2. 达到设定的最大执行时间（如果设置了DURATION参数）
*/

// Configuration - 移到最前面
const WS_IP = __ENV.WS_IP || 'ws://192.168.0.213';
const WS_URL = `${WS_IP}/15109`;
const WS_URL_15101 = `${WS_IP}/15101`;
const ACCOUNT_COUNT = Number(__ENV.ACCOUNT_COUNT) || 2;
const ACCOUNTS = Array.from({ length: ACCOUNT_COUNT }, (_, i) => `laugh_g_${i + 1}`);
const DEBUG = __ENV.DEBUG === 'true'; // Debug flag to control logging
const OUTPUT_FILE = __ENV.OUTPUT_FILE || null; // Output file path for test results
const DURATION = __ENV.DURATION || null; // Duration for test execution (e.g., '10m', '1h')
const PAYOUT_COUNT = Number(__ENV.PAYOUT_COUNT) || 10; // Number of payouts required per account

// Parse GROUP_IDS from environment variable
const GROUP_IDS_STR = __ENV.GROUP_IDS || '';
const GROUP_IDS = GROUP_IDS_STR ? GROUP_IDS_STR.split(',').map(id => parseInt(id.trim())) : [];

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

// Helper function to get random groupID for an account
function getRandomGroupID(account) {
  if (GROUP_IDS.length === 0) {
    return 0; // Default behavior if no groupIDs provided
  }
  const randomIndex = randomIntBetween(0, GROUP_IDS.length - 1);
  const selectedGroupID = GROUP_IDS[randomIndex];
  debugLog(`账号 ${account} 随机选择 groupID: ${selectedGroupID}`);
  return selectedGroupID;
}

export const options = {
  vus: Number(__ENV.ACCOUNT_COUNT) || 2,
  iterations: Number(__ENV.ACCOUNT_COUNT) || 2, // 每个VU只执行一次迭代
  ...(DURATION && { duration: DURATION }), // 如果设置了DURATION，则添加到options中
};

// 添加一个全局变量来跟踪VU完成状态
const vuCompleted = new Set();
// 添加一个全局变量来跟踪每个VU的派彩次数
const vuPayoutCounts = new Map();

// 添加teardown函数来确保所有VU完成后立即停止
export function teardown(data) {
  console.log('所有VU迭代完成，k6测试结束');
}

export default function () {
  const account = ACCOUNTS[__VU % ACCOUNTS.length];
  
  // 检查这个VU是否已经完成
  if (vuCompleted.has(__VU)) {
    console.log(`VU ${__VU} (账号: ${account}) 已经完成，跳过`);
    return;
  }
  
  // 初始化这个VU的派彩计数器
  if (!vuPayoutCounts.has(__VU)) {
    vuPayoutCounts.set(__VU, 0);
  }
  
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
  // console.log(`开始连接15101, 账号: ${account}, sid: ${sid}`);
  
  const res = ws.connect(WS_URL_15101, {}, function (socket) {
    socket.on('open', function () {
      //console.log(`15101 WebSocket连接已打开: ${WS_URL_15101} 账号: ${account}`);
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
      debugLog(`发送15101登录消息: ${account}`);
    });

    let firstTwentyOneFlag = true;
    let modifyBetLimitSendFlag = false;
    let modifyBetLimitFlag = false;
    let recvBetResultFlag = false;
    let betSerialNumber = 1;
    let groupID = getRandomGroupID(account);
    let memberID = 0;
    let pingInterval = null;

    socket.on('message', function (message) {
      try {
        const data = JSON.parse(message);
        if (data.protocol === 0) {
          console.log(`15101 登录成功: ${account}`);
        } else if (data.protocol === 10) {
          // Join table response
          if (data.data && data.data.memberID) {
            memberID = data.data.memberID;
          }
          console.log(`15101 进入桌台成功: ${account}`);
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
              debugLog(`15101 修改限红成功: ${account}`);
          }
        } else if (data.protocol === 25) {
          // Game result
          if (data.data && data.data.groupID !== groupID && groupID !== 0) return;
          debugLog(`15101 得到一局结果: ${account}`);
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
          // console.log(`15101 派彩成功: account: ${account} ${JSON.stringify(message)}`);
          
          // 增加派彩计数
          const currentPayoutCount = vuPayoutCounts.get(__VU) + 1;
          vuPayoutCounts.set(__VU, currentPayoutCount);
          console.log(`账号 ${account} 派彩次数: ${currentPayoutCount}/${PAYOUT_COUNT}`);
          
          // 检查是否达到指定次数的派彩
          if (currentPayoutCount >= PAYOUT_COUNT) {
            // 收到指定次数的派彩信息后，标记VU完成并关闭连接
            console.log(`账号 ${account} 完成${PAYOUT_COUNT}次派彩，准备结束VU迭代`);
            
            // 标记这个VU为已完成
            vuCompleted.add(__VU);
            
            // 清理ping定时器
            if (pingInterval) {
              clearInterval(pingInterval);
            }
            
            // 关闭WebSocket连接
            socket.close();
            
            // 添加一个标记，表示这个VU已经完成
            console.log(`VU迭代完成 - 账号: ${account}`);
          }
        } else if (data.protocol === 38) {
          // Bet time
          const betTimeData = data.data;
          // Only process if the server's groupID matches our randomly selected groupID
          if (betTimeData.groupID !== groupID && groupID !== 0) return;
          if (firstTwentyOneFlag) {
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
            debugLog(`15101 发送登录桌台: groupID: ${groupID} account: ${account}`);
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
    pingInterval = setInterval(function () {
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