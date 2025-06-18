import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  vus: 2, // Number of virtual users
  duration: '10m', // Test duration
};

// Configuration
const WS_IP = __ENV.WS_IP || 'ws://192.168.0.213'; // Get IP from environment variable or use default
const WS_URL = `${WS_IP}/15109`;
const WS_URL_15101 = `${WS_IP}/15101`;
const ACCOUNTS = ['laugh_g_3', 'laugh_g_4'];

export default function () {
  const account = ACCOUNTS[__VU % ACCOUNTS.length];
  
  // Connect to 15109 for authentication
  const authRes = ws.connect(WS_URL, {}, function (socket) {
    socket.on('open', function () {
      console.log(`连接到 ${WS_URL} 账号: ${account}`);
      
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
      const data = JSON.parse(message);
      
      if (data.protocol === 0) {
        // Authentication successful, get SID
        const sid = data.data.sid;
        console.log(`15109登录成功, sid: ${sid}`);
        
        // Connect to 15101 with the SID
        connectTo15101(sid, account);
        
        // Close the authentication socket
        socket.close();
      }
    });

    socket.on('error', function (e) {
      console.error(`15109连接错误 ${account}:`, e);
    });

    socket.on('close', function () {
      console.log(`15109连接关闭 ${account}`);
    });
  });

  check(authRes, { '15109连接成功': (r) => r && r.status === 101 });
}

function connectTo15101(sid, account) {
  const res = ws.connect(WS_URL_15101, {}, function (socket) {
    socket.on('open', function () {
      console.log(`连接到 ${WS_URL_15101} 账号: ${account}`);
      
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
    });

    let firstTwentyOneFlag = true;
    let modifyBetLimitFlag = false;
    let recvBetResultFlag = false;
    let betSerialNumber = 1;
    let groupID = 0;

    socket.on('message', function (message) {
      const data = JSON.parse(message);
      
      if (data.protocol === 1) {
        console.log(`15101 登录成功: ${JSON.stringify(message)}`);
      } else if (data.protocol === 10) {
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
        console.log(`15101 修改限红`);
      } else if (data.protocol === 60) {
        console.log(`15101 修改限红成功: ${JSON.stringify(message)}`);
        modifyBetLimitFlag = true;
      } else if (data.protocol === 25) {
        const gameResultData = data.data;
        if (gameResultData.groupID !== groupID && groupID !== 0) return;
        console.log(`15101 得到一局结果: ${JSON.stringify(message)}`);
        if (modifyBetLimitFlag) {
          recvBetResultFlag = true;
        }
      } else if (data.protocol === 38) {
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
          console.log(`15101 join table success`);
          console.log(`15101 得到下注时间剩余秒数信息第一次: ${JSON.stringify(betTimeData)}`);
          return;
        }

        console.log(`15101 得到下注时间剩余秒数: ${JSON.stringify(betTimeData)}`);
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
          console.log('15101 投注成功: ' + bettingMsg);
          recvBetResultFlag = false;
          betSerialNumber++;
        }
      } else if (data.protocol === 22) {
        // Handle betting response
        if (data.data.groupID === groupID) {
          console.log(`15101 投注成功回复: ${JSON.stringify(message)}`);
        }
      } else if (data.protocol === 31) {
        console.log(`15101 派彩成功: ${JSON.stringify(message)}`);
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

  check(res, { '15101连接成功': (r) => r && r.status === 101 });
} 