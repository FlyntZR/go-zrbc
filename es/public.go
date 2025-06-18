package es

import (
	"context"
	"encoding/json"
	"go-zrbc/db"
	"go-zrbc/pkg/xlog"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/shopspring/decimal"
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

// GetBet01ListForUnsettledReportEs queries unsettled bet01 data from Elasticsearch
func (c *Client) GetBet01ListForUnsettledReportEs(ctx context.Context, date time.Time) ([]*db.Bet01Summary, error) {
	boolQuery := elastic.NewBoolQuery()

	// Filter by bet07 (accounting date)
	dateStr := date.Format("2006-01-02T15:04:05Z")
	boolQuery.Must(elastic.NewTermQuery("bet07", dateStr))

	// Exclude bet02 == 301
	boolQuery.MustNot(elastic.NewTermQuery("bet02", 301))

	// Exclude bet30 == "Y" (cancelled)
	boolQuery.Must(elastic.NewTermQuery("bet30", 2))

	boolQuery.Must(elastic.NewTermQuery("settle", false))

	// Exclude bet01 that exist in bet02 (simulate: bet01 NOT IN (SELECT bet01 FROM bet02))
	// This is tricky in ES; assume a field like 'is_settled' or similar, or skip this filter if not available.
	// If ES has a field 'is_settled' or similar, use it. Otherwise, this filter may need to be handled at index time.
	// boolQuery.Must(elastic.NewTermQuery("is_settled", false)) // Uncomment if such a field exists

	searchResult, err := c.client.Search().
		Index("bet01_report_index").
		Query(boolQuery).
		Size(10000).
		Do(ctx)
	if err != nil {
		xlog.Errorf("Elasticsearch query failed: %v, query: %+v", err, boolQuery)
		return nil, err
	}

	xlog.Debugf("debug to get bet01 list from ES, searchResult: %v", searchResult)

	var results []*db.Bet01Summary
	for _, hit := range searchResult.Hits.Hits {
		var src map[string]interface{}
		if err := json.Unmarshal(hit.Source, &src); err != nil {
			xlog.Errorf("Failed to unmarshal hit: %v", err)
			continue
		}
		item := &db.Bet01Summary{}
		if v, ok := src["betId"].(float64); ok {
			item.BetID = int64(v)
		}
		if v, ok := src["gid"].(float64); ok {
			item.GID = int(v)
		}
		if v, ok := src["event"].(string); ok {
			item.Event, _ = decimal.NewFromString(v)
		} else if v, ok := src["event"].(float64); ok {
			item.Event = decimal.NewFromFloat(v)
		}
		if v, ok := src["eventChild"].(float64); ok {
			item.EventChild = int(v)
		}
		if v, ok := src["userId"].(float64); ok {
			item.ID = int(v)
		}
		if v, ok := src["betTime"].(string); ok {
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				item.BetTime = t
			}
		} else if v, ok := src["betTime"].(float64); ok {
			item.BetTime = time.Unix(int64(v)/1000, 0)
		}
		if v, ok := src["betResult"].(string); ok {
			item.BetResult = v
		}
		if v, ok := src["bet"].(string); ok {
			item.Bet, _ = decimal.NewFromString(v)
		} else if v, ok := src["bet"].(float64); ok {
			item.Bet = decimal.NewFromFloat(v)
		}
		if v, ok := src["aid"].(float64); ok {
			item.AID = int(v)
		}
		if v, ok := src["tableId"].(float64); ok {
			item.TableID = int(v)
		}
		if v, ok := src["commission"].(float64); ok {
			item.Commission = int(v)
		}
		if v, ok := src["event"].(string); ok {
			item.Round, _ = decimal.NewFromString(v)
		} else if v, ok := src["round"].(float64); ok {
			item.Round = decimal.NewFromFloat(v)
		}
		if v, ok := src["eventChild"].(float64); ok {
			item.SubRound = int(v)
		}
		// gname and user may be indexed in ES, or need enrichment
		if v, ok := src["gameName"].(string); ok {
			item.GName = v
		}
		if v, ok := src["username"].(string); ok {
			item.User = v
		}
		results = append(results, item)
	}
	xlog.Debugf("debug to get bet01 list from ES, results: %v", results)
	return results, nil
}
