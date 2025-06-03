package mq

import (
	"go-zrbc/pkg/xlog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	TopicBibSwapOrder = "bib_push_go_swap_order"
	// TopicBibSpotOrder = "bib_push_go_spot_order"
	TopicBibSpotOrder     = "spot_core_cmd_result"
	TopicBibSpotTicket    = "spot_ticket"
	TopicBibSpotKline     = "spot_kline"
	TopicBibDepth         = "bib_push_go_depth"
	TopicBibKline         = "BIB_PUSH_GO_KLINE"
	TopicRocketmqBibTrade = "SPOT_SELF_TRADE"
	TopicDBToolKline      = "write_kline_topic" // uat rocketmq地址配置：10.8.34.64:9876;10.8.34.65:9876;10.8.34.66:9876
	//TopicClientKline   = "quotation"
	// 本机测试用
	TopicClientKline = "BIB_PUSH_GO_KLINE"
	TopicStatsReport = "bib_push_go_stats_report"
	TopicBibApi      = "bib_push_go_api"

	UniqueGroupID = "ws_channel_1"
	SharedGroupID = "ws_channel"
	DepthGroupID  = "depth_group_id"

	// Time to auto commit.
	commitInterval = 5 * time.Second

	// print log when commit how many times, is 2 min.
	logFrequency = 60
)

func NewWriter(addrs []string, topic string) *kafka.Writer {
	w := &kafka.Writer{
		BatchTimeout: time.Millisecond * 100,
		Addr:         kafka.TCP(addrs...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Async:        true, // make the writer asynchronous
		Completion: func(messages []kafka.Message, err error) {
		},
	}
	return w
}

type TestHash struct {
}

func (h *TestHash) Balance(msg kafka.Message, partitions ...int) int {
	return msg.Partition
}

func NewWriterWithBalance(addrs []string, topic string) *kafka.Writer {
	w := &kafka.Writer{
		BatchTimeout: time.Millisecond * 100,
		Addr:         kafka.TCP(addrs...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Async:        true, // make the writer asynchronous
		Completion: func(messages []kafka.Message, err error) {
		},
	}
	w.Balancer = &TestHash{}
	return w
}

func NewReaderInPartition(addrs []string, topic string) *kafka.Reader {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   addrs,
		Topic:     topic,
		Partition: 0,
		MinBytes:  1, // 10KB
		MaxWait:   time.Millisecond * 500,
		MaxBytes:  10e6, // 10MB
	})
	//r.SetOffsetAt(context.TODO(), time.Now())
	//r.SetOffset(100000)
	return r
}

func NewReaderInPartitionNum(addrs []string, topic string, partitionNum int) *kafka.Reader {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   addrs,
		Topic:     topic,
		Partition: partitionNum,
		MinBytes:  1, // 10KB
		MaxWait:   time.Millisecond * 500,
		MaxBytes:  10e6, // 10MB
	})
	//r.SetOffsetAt(context.TODO(), time.Now())
	//r.SetOffset(100000)
	return r
}

func NewReaderInGroup(addrs []string, topic string, groupID string) *kafka.Reader {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1, // 10KB
		MaxWait:        time.Millisecond * 500,
		MaxBytes:       10e6, // 10MB
		CommitInterval: 2 * time.Second,
		StartOffset:    kafka.LastOffset,
	})
	return r
}

func NewConsumerGroup(addrs []string, topic string, groupID string) (*kafka.ConsumerGroup, error) {
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      groupID,
		Brokers: addrs,
		Topics:  []string{topic},
	})
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	return group, nil
}

type CommitInfo struct {
	sync.RWMutex
	Topic      string            `json:"topic"`
	Generation *kafka.Generation `json:"generation"`
	Partition  int               `json:"partition"`
	Offset     int64             `json:"offset"`
}

func CommitKafkaOffset(commitInfo <-chan *CommitInfo) {
	ticker := time.NewTicker(commitInterval)
	defer ticker.Stop()
	var c *CommitInfo
	logInterval := 0
	for {
		select {
		case <-ticker.C:
			if c != nil {
				c.RLock()
				defer c.RUnlock()
				if err := c.Generation.CommitOffsets(map[string]map[int]int64{c.Topic: {c.Partition: c.Offset + 1}}); err != nil {
					xlog.Errorf("commit offset err, commit info:(%+v), err:(%+v)\n", c, err)
				}
				logInterval = logInterval + 1
				if logInterval == logFrequency {
					// print commit log per 2 min.
					logInterval = 0
					xlog.Infof("commit offset :(%+v)\n", c)
				}
			}
		case c = <-commitInfo:
		}
	}
}
