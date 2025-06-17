package es

import (
	"context"
	"encoding/json"
	"go-zrbc/pkg/xlog"
	"time"

	"github.com/olivere/elastic/v7"
)

// GetBet02ListForDateTimeReportEs queries bet02 data from Elasticsearch
func (c *Client) GetBet02ListForDateTimeReportEs(ctx context.Context, memberID int64, agentID int64, startTime, endTime int64, dataType, timeType int, gameNo1, gameNo2 string) ([]map[string]interface{}, error) {
	// Create bool query
	boolQuery := elastic.NewBoolQuery()

	// Add member/agent filter
	if memberID != 0 {
		boolQuery.Must(elastic.NewTermQuery("bet05", float64(memberID)))
	} else if agentID != 0 {
		boolQuery.Must(elastic.NewTermQuery("bet22", float64(agentID)))
	}

	// Add data type filter
	if dataType == 0 {
		boolQuery.MustNot(elastic.NewMatchQuery("bet09", "Tip_1_"))
	} else if dataType == 1 {
		boolQuery.Must(elastic.NewMatchQuery("bet09", "Tip_1_"))
	}

	// Add game number filters
	if gameNo1 != "" {
		boolQuery.Must(elastic.NewTermQuery("bet03", gameNo1))
		if gameNo2 != "" {
			boolQuery.Must(elastic.NewTermQuery("bet04", gameNo2))
		}
	}

	// Add time range filters using proper date format
	startTimeStr := time.Unix(startTime, 0).Format("2006-01-02T15:04:05Z")
	endTimeStr := time.Unix(endTime, 0).Format("2006-01-02T15:04:05Z")
	if timeType == 1 {
		boolQuery.Must(elastic.NewRangeQuery("updatetime").Gte(startTimeStr).Lte(endTimeStr))
	} else {
		boolQuery.Must(elastic.NewRangeQuery("bet08").Gte(startTimeStr).Lte(endTimeStr))
	}

	// Execute search with reasonable size limit and timeout
	searchResult, err := c.client.Search().
		Index("bet02_report_index").
		Query(boolQuery).
		Size(10000).    // Increased size limit for debugging
		Timeout("60s"). // Increased timeout
		Do(ctx)
	if err != nil {
		xlog.Errorf("Elasticsearch query failed: %v, query: %+v", err, boolQuery)
		return nil, err
	}

	// Process results
	var results []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		xlog.Infof("hit.Source: %+v", hit.Source)
		var bet02 map[string]interface{}
		if err := json.Unmarshal(hit.Source, &bet02); err != nil {
			xlog.Errorf("Failed to unmarshal hit: %v", err)
			continue
		}
		results = append(results, bet02)
	}

	return results, nil
}

// GetInOutMsEs queries in_out_m data from Elasticsearch
func (c *Client) GetInOutMsEs(ctx context.Context, mIDs []int64, orderID, order string, startTime, endTime int64) ([]map[string]interface{}, error) {
	// Create bool query
	boolQuery := elastic.NewBoolQuery()

	// Add member IDs filter
	if len(mIDs) > 0 {
		xlog.Infof("mIDs: %v", mIDs)
		terms := make([]interface{}, len(mIDs))
		for i, v := range mIDs {
			terms[i] = v
		}
		boolQuery.Must(elastic.NewTermsQuery("iom003", terms...))
	}

	// Add operation code filter
	boolQuery.Must(elastic.NewTermsQuery("iom005", "121", "122", "501", "502", "504"))

	// Add order ID filter
	if orderID != "" {
		xlog.Infof("orderID: %v", orderID)
		boolQuery.Must(elastic.NewTermQuery("iom001", orderID))
	}

	// Add order filter
	if order != "" {
		xlog.Infof("order: %v", order)
		boolQuery.Must(elastic.NewTermQuery("iom008", order))
	}

	// Add time range filters if no specific order is provided
	if orderID == "" && order == "" {
		var startTimeStr, endTimeStr string

		if startTime == 0 || endTime == 0 {
			// Default to last hour if no time range specified
			now := time.Now()
			startTimeStr = now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
			endTimeStr = now.Format("2006-01-02 15:04:05")
		} else {
			startTimeStr = time.Unix(startTime, 0).Format("2006-01-02 15:04:05")
			endTimeStr = time.Unix(endTime, 0).Format("2006-01-02 15:04:05")
		}

		boolQuery.Must(elastic.NewRangeQuery("iom002").Gte(startTimeStr).Lte(endTimeStr))
	}

	// Execute search
	searchResult, err := c.client.Search().
		Index("in_out_m_index").
		Query(boolQuery).
		Do(ctx)
	if err != nil {
		xlog.Errorf("error to db transaction, err: %v", err)
		return nil, err
	}

	// Process results
	var results []map[string]interface{}
	for _, hit := range searchResult.Hits.Hits {
		xlog.Infof("hit: %v", hit)
		var inOutM map[string]interface{}
		if err := json.Unmarshal(hit.Source, &inOutM); err != nil {
			continue
		}
		results = append(results, inOutM)
	}

	return results, nil
}

// GetBet02ListForReportDetailEs queries bet02 detail from Elasticsearch by betID
func (c *Client) GetBet02ListForReportDetailEs(ctx context.Context, betID int64) (map[string]interface{}, error) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewTermQuery("bet01", float64(betID)))

	searchResult, err := c.client.Search().
		Index("bet02_report_index").
		Query(boolQuery).
		Size(1).
		Do(ctx)
	if err != nil {
		xlog.Errorf("Elasticsearch query failed: %v", err)
		return nil, err
	}

	if searchResult.TotalHits() == 0 {
		return nil, nil
	}

	hit := searchResult.Hits.Hits[0]
	var bet02 map[string]interface{}
	if err := json.Unmarshal(hit.Source, &bet02); err != nil {
		xlog.Errorf("Failed to unmarshal hit: %v", err)
		return nil, err
	}

	return bet02, nil
}
