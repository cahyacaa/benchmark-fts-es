package main

import (
	"fmt"
	"log"

	"github.com/bencmark-fts-es/database"
)

func main() {
	// PostgreSQL Setup
	pgConnStr := "postgres://postgres:password@localhost:5432/benchmarkdb?sslmode=disable"
	pg := database.NewPostgresSearch(pgConnStr)

	if err := pg.CreateTable(); err != nil {
		log.Fatalf("PostgreSQL table creation failed: %v", err)
	}
	if err := pg.CreateFullTextIndex(); err != nil {
		log.Fatalf("PostgreSQL index creation failed: %v", err)
	}

	// Elasticsearch Setup
	es := database.NewElasticsearchSearch([]string{"http://localhost:9200"})
	if err := es.CreateIndex(); err != nil {
		log.Fatalf("Elasticsearch index creation failed: %v", err)
	}

	// Sample documents
	documents := []struct {
		title   string
		content string
	}{
		{"Tech Innovation", "Exploring cutting-edge technologies in artificial intelligence"},
		{"Machine Learning", "Deep learning algorithms are revolutionizing data science"},
		// Add more sample documents
	}

	// Index documents
	for _, doc := range documents {
		if err := pg.IndexDocument(doc.content, doc.title); err != nil {
			log.Printf("PostgreSQL indexing error: %v", err)
		}

		if err := es.IndexDocument(doc.content, doc.title); err != nil {
			log.Printf("Elasticsearch indexing error: %v", err)
		}
	}

	// Benchmark Searches
	benchmarkSearches(pg, es)
}

func benchmarkSearches(pg *database.PostgresSearch, es *database.ElasticsearchSearch) {
	searches := []string{
		"technology",
		"machine learning",
		"artificial intelligence",
	}

	for _, query := range searches {
		fmt.Printf("Searching for: %s\n", query)

		// PostgreSQL Search
		pgResults, pgDuration, err := pg.Search(query, 10)
		if err != nil {
			log.Printf("PostgreSQL search error: %v", err)
			continue
		}

		// Elasticsearch Search
		esResults, esDuration, err := es.Search(query, 10)
		if err != nil {
			log.Printf("Elasticsearch search error: %v", err)
			continue
		}

		// Print Results and Performance
		fmt.Printf("PostgreSQL Results (Duration: %v):\n", pgDuration)
		for _, r := range pgResults {
			fmt.Println(r)
		}

		fmt.Printf("\nElasticsearch Results (Duration: %v):\n", esDuration)
		for _, r := range esResults {
			fmt.Println(r)
		}

		fmt.Println("---")
	}
}
