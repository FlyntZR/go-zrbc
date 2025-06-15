package es

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

type Client struct {
	client *elastic.Client
}

func NewClient(url string) (*Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(128*time.Millisecond, 513*time.Millisecond))),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %v", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Search(ctx context.Context, index string, query elastic.Query, from, size int) (*elastic.SearchResult, error) {
	searchResult, err := c.client.Search().
		Index(index).
		Query(query).
		From(from).
		Size(size).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %v", err)
	}

	return searchResult, nil
}

func (c *Client) Index(ctx context.Context, index string, id string, doc interface{}) error {
	_, err := c.client.Index().
		Index(index).
		Id(id).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to index document: %v", err)
	}

	return nil
}

func (c *Client) BulkIndex(ctx context.Context, index string, docs []interface{}) error {
	bulk := c.client.Bulk()
	for _, doc := range docs {
		req := elastic.NewBulkIndexRequest().
			Index(index).
			Doc(doc)
		bulk.Add(req)
	}

	_, err := bulk.Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to bulk index documents: %v", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c.client != nil {
		c.client.Stop()
	}
	return nil
}
