package es

import (
	"context"
	"encoding/json"
	"time"

	"github.com/olivere/elastic/v7"
)

// GetBet02ListForMemberReportEs queries bet02 data from Elasticsearch
func (c *Client) GetBet02ListForMemberReportEs(ctx context.Context, memberID int64, agentID int64, startTime, endTime int64, dataType, timeType int, gameNo1, gameNo2 string) ([]map[string]interface{}, error) {
	// Create bool query
	boolQuery := elastic.NewBoolQuery()

	// Add member/agent filter
	if memberID != 0 {
		boolQuery.Must(elastic.NewTermQuery("bet05", memberID))
	} else {
		boolQuery.Must(elastic.NewTermQuery("bet22", agentID))
	}

	// Add data type filter
	if dataType == 0 {
		boolQuery.MustNot(elastic.NewTermQuery("bet09", "Tip_1_"))
	} else if dataType == 1 {
		boolQuery.Must(elastic.NewTermQuery("bet09", "Tip_1_"))
	}

	// Add game number filters
	if gameNo1 != "" {
		boolQuery.Must(elastic.NewTermQuery("bet03", gameNo1))
		if gameNo2 != "" {
			boolQuery.Must(elastic.NewTermQuery("bet04", gameNo2))
		}
	}

	// Add time range filters
	startTimeStr := time.Unix(startTime, 0).Format("2006-01-02 15:04:05")
	endTimeStr := time.Unix(endTime, 0).Format("2006-01-02 15:04:05")
	if timeType == 1 {
		boolQuery.Must(elastic.NewRangeQuery("updatetime").Gte(startTimeStr).Lte(endTimeStr))
	} else {
		boolQuery.Must(elastic.NewRangeQuery("bet08").Gte(startTimeStr).Lte(endTimeStr))
	}

	// Execute search
	searchResult, err := c.client.Search().
		Index("bet02").
		Query(boolQuery).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// Process results
	var results []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		var bet02 map[string]interface{}
		if err := json.Unmarshal(hit.Source, &bet02); err != nil {
			continue
		}
		results = append(results, bet02)
	}

	return results, nil
}
