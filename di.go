package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codingconcepts/di/pkg/flags"
	"github.com/codingconcepts/di/pkg/model"
	"github.com/codingconcepts/di/pkg/runner"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

func main() {
	log.SetFlags(0)

	file := flag.String("file", "", "absolute or relative path to the CSV file to import")
	url := flag.String("url", "", "database connection string")
	schema := flag.String("schema", "public", "schema in which your table resides")
	table := flag.String("table", "", "name of the table to import into")
	batchSize := flag.Int("batch", 10000, "import batch size")

	var formatHelpers flags.StringSlice
	flag.Var(&formatHelpers, "fmt", "format helper (in the form of column_name:format)")

	flag.Parse()

	*url = "postgres://root@localhost:26257/store?sslmode=disable"
	*file = "examples/simple/csvs/customer.csv"

	if *file == "" || *url == "" {
		flag.Usage()
		os.Exit(2)
	}

	db, err := pgxpool.New(context.Background(), *url)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = db.Ping(timeout); err != nil {
		log.Fatalf("error testing database connection: %v", err)
	}

	if *table == "" {
		base := filepath.Base(*file)
		*table = strings.TrimSuffix(base, filepath.Ext(base))
	}

	types, err := fetchTableInformation(db, *schema, *table)
	if err != nil {
		log.Fatalf("error fetching table information: %v", err)
	}

	addFormatters(types, formatHelpers)

	runner := runner.New(db, *table, types, *batchSize)

	csv, err := os.Open(*file)
	if err != nil {
		log.Fatalf("error opening csv file for reading: %v", err)
	}
	defer csv.Close()

	if err = runner.StreamCSV(csv); err != nil {
		log.Fatalf("error streaming csv to database: %v", err)
	}
}

func fetchTableInformation(db *pgxpool.Pool, schema, table string) (model.ColumnTypes, error) {
	const stmt = `SELECT ordinal_position, column_name, udt_name, is_nullable
								FROM information_schema.columns
								WHERE table_name = $1
								AND table_schema = $2
								ORDER BY ordinal_position`

	rows, err := db.Query(context.Background(), stmt, table, schema)
	if err != nil {
		return nil, fmt.Errorf("fetching column information: %w", err)
	}

	log.Printf("Table: %s", table)
	log.Printf("Columns: %s\n\n", table)
	defer log.Println()

	types := model.ColumnTypes{}
	for rows.Next() {
		var c model.Column
		var rawNullable string

		if err = rows.Scan(&c.Ordinal, &c.Name, &c.Type, &rawNullable); err != nil {
			return nil, fmt.Errorf("scanning column information: %w", err)
		}

		c.Nullable = strings.EqualFold(rawNullable, "YES")

		types[c.Name] = &c

		log.Printf(" %d. %s (%s)%s", c.Ordinal, c.Name, c.Type, lo.Ternary(c.Nullable, "- NULL", ""))
	}

	if len(types) == 0 {
		return nil, fmt.Errorf("no types found, check database name")
	}

	return types, nil
}

func addFormatters(columns model.ColumnTypes, formatters flags.StringSlice) error {
	for _, f := range formatters {
		parts := strings.Split(f, ":")
		name := parts[0]
		format := parts[1]

		c, ok := columns[name]
		if !ok {
			return fmt.Errorf("missing column for formatter: %q", f)
		}

		c.Format = format
	}

	return nil
}
