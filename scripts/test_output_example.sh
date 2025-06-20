#!/bin/bash

# 测试脚本示例：使用k6的正确输出参数

echo "=== k6测试结果输出到文件示例 ==="

# 示例1：输出到JSON文件 (正确的用法)
echo "1. 输出到JSON文件 (使用 --out json=filename):"
k6 run --out json=test_results.json \
       -e WS_IP=wss://a45gs-t.wmetg.com \
       -e ACCOUNT_COUNT=2 \
       -e GROUP_IDS=6,7,8 \
       -e DEBUG=true \
       scripts/betting_test_new.js

echo "测试结果已保存到 test_results.json"

# 示例2：输出到CSV文件
echo "2. 输出到CSV文件:"
k6 run --out csv=test_results.csv \
       -e WS_IP=wss://a45gs-t.wmetg.com \
       -e ACCOUNT_COUNT=2 \
       -e GROUP_IDS=6,7,8 \
       -e DEBUG=true \
       scripts/betting_test_new.js

echo "测试结果已保存到 test_results.csv"

# 示例3：同时输出到多个格式
echo "3. 同时输出到多个格式:"
k6 run --out json=test_results.json \
       --out csv=test_results.csv \
       --out influxdb=http://localhost:8086/k6 \
       -e WS_IP=wss://a45gs-t.wmetg.com \
       -e ACCOUNT_COUNT=2 \
       -e GROUP_IDS=6,7,8 \
       -e DEBUG=true \
       scripts/betting_test_new.js

echo "=== 正确的k6输出格式 ==="
echo "JSON格式: --out json=results.json"
echo "CSV格式: --out csv=results.csv"
echo "InfluxDB: --out influxdb=http://localhost:8086/k6"
echo "Prometheus: --out prometheus=localhost:9090/k6"
echo "Kafka: --out kafka=localhost:9092"
echo "CloudWatch: --out cloudwatch=us-east-1"

echo ""
echo "=== 你的原始命令应该修改为 ==="
echo "k6 run --out json=test_results.json \\"
echo "       -e WS_IP=wss://a45gs-t.wmetg.com \\"
echo "       -e ACCOUNT_COUNT=20 \\"
echo "       -e DEBUG=true \\"
echo "       -e GROUP_IDS=6,7,8 \\"
echo "       scripts/betting_test_new.js" 