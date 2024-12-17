package database

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticsearchSearch struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchSearch(addresses []string) *ElasticsearchSearch {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}
	return &ElasticsearchSearch{
		client: client,
		index:  "documents",
	}
}

func (e *ElasticsearchSearch) CreateIndex() error {
	req := esapi.IndicesCreateRequest{
		Index: e.index,
		Body: bytes.NewReader([]byte(`{
            "mappings": {
                "properties": {
                    "content": { "type": "text" },
                    "title": { "type": "text" }
                }
            }
        }`)),
	}
	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (e *ElasticsearchSearch) IndexDocument(content, title string) error {
	doc := map[string]interface{}{
		"content": content,
		"title":   title,
	}

	req := esapi.IndexRequest{
		Index:      e.index,
		DocumentID: "",
		Body:       bytes.NewReader(jsonMarshal(doc)),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (e *ElasticsearchSearch) Search(query string, limit int) ([]string, time.Duration, error) {
	start := time.Now()

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"content", "title"},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, 0, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(e.index),
		e.client.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	var results []string
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		results = append(results, source["title"].(string))
	}

	duration := time.Since(start)
	return results, duration, nil
}

func jsonMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
