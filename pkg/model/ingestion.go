package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Ingester defines the behaviour of something that can ingest data
// into the database.
type Ingester interface {
	Ingest(db *pgxpool.Pool, table string, header []string, rows [][]any) error
}

// CopyFromIngester is an implementation of Ingester that ingests data
// by way of the COPY FROM mechanism.
type CopyFromIngester struct{}

func (i *CopyFromIngester) Ingest(db *pgxpool.Pool, table string, header []string, rows [][]any) error {
	_, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{table},
		header,
		pgx.CopyFromRows(rows),
	)

	return err
}

// BatchInsertIngester is an implementation of Ingester that ingests data
// by way of the pgx.Batch mechanism.
type BatchInsertIngester struct{}

func (i *BatchInsertIngester) Ingest(db *pgxpool.Pool, table string, header []string, rows [][]any) error {
	var batch pgx.Batch

	stmt := insertStatement(table, header)
	for _, row := range rows {
		batch.Queue(stmt, row...)
	}

	_, err := db.SendBatch(context.Background(), &batch).Exec()
	return err
}

func insertStatement(table string, header []string) string {
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(header, ","),
		headerDollars(header),
	)
}

func headerDollars(slice []string) string {
	var parts []string
	for i := range slice {
		parts = append(parts, fmt.Sprintf("$%d", i+1))
	}
	return strings.Join(parts, ",")
}
