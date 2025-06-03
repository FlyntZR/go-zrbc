# Flink Doris Integration Tests

This directory contains integration tests for Apache Flink and Apache Doris.

## Prerequisites

Before running the tests, make sure you have the following components set up:

1. Apache Flink (version 1.15.0 or later)
   - Running Flink cluster with JobManager accessible at `localhost:8081`
   - Flink SQL client available

2. Apache Doris (version 1.2.0 or later)
   - Running Doris cluster with FE accessible at `localhost:9030`
   - Doris BE accessible at `localhost:8030`

3. Required Go dependencies:
   ```bash
   go get github.com/apache/doris-go/doris-connector-for-go
   go get github.com/apache/flink-go/client
   go get github.com/stretchr/testify/assert
   ```

## Configuration

The test uses the following default configuration:

- Doris connection:
  - Host: localhost
  - Port: 9030 (FE), 8030 (BE)
  - User: root
  - Password: (empty)
  - Database: test

- Flink connection:
  - Host: localhost
  - Port: 8081

To modify these settings, update the constants in `flink_doris_test.go`.

## Running the Tests

To run all tests:
```bash
go test -v ./...
```

To run specific Flink Doris integration test:
```bash
go test -v -run TestFlinkDorisIntegration
```

## Test Description

The test case `TestFlinkDorisIntegration` performs the following:

1. Establishes connections to both Flink and Doris
2. Creates a test table in Doris
3. Submits a Flink SQL job that:
   - Creates a source table using Flink's datagen connector
   - Creates a sink table connected to Doris
   - Inserts data from source to sink
4. Verifies that data is successfully written to Doris
5. Cleans up by dropping the test table and canceling the Flink job

## Troubleshooting

1. If the test fails to connect to Doris:
   - Verify Doris FE is running: `curl http://localhost:9030`
   - Check Doris user permissions

2. If the test fails to connect to Flink:
   - Verify Flink cluster is running: `curl http://localhost:8081`
   - Check Flink logs for any errors

3. If no data is written to Doris:
   - Increase the wait time (currently 30 seconds)
   - Check Flink job status in the Flink UI
   - Verify Doris BE is running and accessible 