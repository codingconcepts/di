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

var version string

func main() {
	log.SetFlags(0)

	file := flag.String("file", "", "absolute or relative path to the CSV file to import")
	url := flag.String("url", "", "database connection string")
	schema := flag.String("schema", "public", "schema in which your table resides")
	table := flag.String("table", "", "name of the table to import into")
	batchSize := flag.Int("batch", 10000, "import batch size")
	displayVersion := flag.Bool("version", false, "display version information")

	formatHelpers := flags.NewStringSlice()
	flag.Var(formatHelpers, "fmt", "format helper (in the form of data_type:format)")

	flag.Parse()

	if *displayVersion {
		log.Println(version)
		os.Exit(0)
	}

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

	runner := runner.New(db, *table, types, *batchSize, formatHelpers)

	csv, err := os.Open(*file)
	if err != nil {
		log.Fatalf("error opening csv file for reading: %v", err)
	}
	defer func() {
		if err = csv.Close(); err != nil {
			log.Fatalf("error closing csv file %v", err)
		}
	}()

	if err = runner.StreamCSV(csv); err != nil {
		log.Fatalf("error streaming csv to database: %v", err)
	}

	log.Println("\n Finished")
}

func fetchTableInformation(db *pgxpool.Pool, schema, table string) (model.ColumnTypes, error) {
	const stmt = `SELECT ordinal_position, column_name, udt_name, is_nullable, is_generated
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
		var rawGenerated string

		if err = rows.Scan(&c.Ordinal, &c.Name, &c.Type, &rawNullable, &rawGenerated); err != nil {
			return nil, fmt.Errorf("scanning column information: %w", err)
		}

		c.Nullable = strings.EqualFold(rawNullable, "YES")
		c.IsGenerated = strings.EqualFold(rawGenerated, "ALWAYS")

		types[c.Name] = &c

		log.Printf(
			" %d. %s (%s)%s%s",
			c.Ordinal,
			c.Name,
			c.Type,
			lo.Ternary(c.Nullable, " - NULL", ""),
			lo.Ternary(c.IsGenerated, " - GENERATED", ""),
		)
	}

	if len(types) == 0 {
		return nil, fmt.Errorf("no types found, check database name")
	}

	return types, nil
}
