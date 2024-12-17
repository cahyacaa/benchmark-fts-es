package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/lib/pq"
)

type PostgresSearch struct {
    db *sql.DB
}

func NewPostgresSearch(connectionString string) *PostgresSearch {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }
    return &PostgresSearch{db: db}
}

func (p *PostgresSearch) CreateTable() error {
    _, err := p.db.Exec(`
        CREATE TABLE IF NOT EXISTS documents (
            id SERIAL PRIMARY KEY,
            content TEXT,
            title TEXT,
            search_vector tsvector
        )
    `)
    return err
}

func (p *PostgresSearch) CreateFullTextIndex() error {
    _, err := p.db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_fts 
        ON documents USING GIN(search_vector)
    `)
    return err
}

func (p *PostgresSearch) IndexDocument(content, title string) error {
    _, err := p.db.Exec(`
        INSERT INTO documents (content, title, search_vector) 
        VALUES ($1, $2, 
            to_tsvector('english', $1 || ' ' || $2)
        )
    `, content, title)
    return err
}

func (p *PostgresSearch) Search(query string, limit int) ([]string, time.Duration, error) {
    start := time.Now()
    rows, err := p.db.Query(`
        SELECT title 
        FROM documents 
        WHERE search_vector @@ to_tsquery('english', $1)
        LIMIT $2
    `, query, limit)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var results []string
    for rows.Next() {
        var title string
        if err := rows.Scan(&title); err != nil {
            return nil, 0, err
        }
        results = append(results, title)
    }

    duration := time.Since(start)
    return results, duration, nil
}