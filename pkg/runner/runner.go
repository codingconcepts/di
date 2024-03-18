package runner

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"github.com/codingconcepts/di/pkg/flags"
	"github.com/codingconcepts/di/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

// Runner holds the runtime properties of the application.
type Runner struct {
	db            *pgxpool.Pool
	table         string
	types         model.ColumnTypes
	formatHelpers map[string]string
	batchSize     int
}

func New(db *pgxpool.Pool, table string, types model.ColumnTypes, batchSize int, formatHelpers *flags.StringSlice) *Runner {
	return &Runner{
		db:            db,
		table:         table,
		types:         types,
		batchSize:     batchSize,
		formatHelpers: formatHelpers.ToMap(),
	}
}

func (runner *Runner) StreamCSV(r io.ReadSeeker) error {
	// Read file lines for progress tracking.
	lines, err := fileLines(r)
	if err != nil {
		return fmt.Errorf("reading file lines: %w", err)
	}

	if _, err = r.Seek(0, 0); err != nil {
		return fmt.Errorf("error resetting file reader: %w", err)
	}

	csvReader := csv.NewReader(r)

	// Assume first row is header.
	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("reading csv header: %w", err)
	}

	// Validate header against table columns.
	if err = validateColumns(header, runner.types); err != nil {
		return fmt.Errorf("validating columns: %w", err)
	}

	// Determine ingestion mechanism based on whether there are generated columns
	// without explicit values in the csv file.
	ingester := createIngester(header, runner.types)
	fmt.Printf("Ingestion mechanism: %T\n", ingester)

	i := 1
	rows := [][]any{}
	for {
		record, err := csvReader.Read()

		// Break if we've reached the end of the file (don't return, as
		// there might be rows to flush).
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading csv file: %w", err)
		}

		args, err := runner.csvLineToArgs(header, record)
		if err != nil {
			return fmt.Errorf("converting csv line to argsL %w", err)
		}
		rows = append(rows, args)

		if i%runner.batchSize == 0 {
			if err = ingester.Ingest(runner.db, runner.table, header, rows); err != nil {
				return fmt.Errorf("flushing rows: %w", err)
			}

			log.Printf("%d/%d rows copied", i, lines)
			rows = [][]any{}
		}

		i++
	}

	if len(rows) > 0 {
		if err = ingester.Ingest(runner.db, runner.table, header, rows); err != nil {
			return fmt.Errorf("flushing rows: %w", err)
		}
	}

	return nil
}

func validateColumns(headers []string, types model.ColumnTypes) error {
	for name, c := range types {
		// Find header and if missing (and type is non-nullable) error.
		_, ok := lo.Find(headers, func(h string) bool {
			return h == name
		})

		if !ok && !c.Nullable && !c.IsGenerated && !c.HasDefault {
			return fmt.Errorf("missing non-nullable column %q", name)
		}
	}

	return nil
}

func createIngester(headers []string, types model.ColumnTypes) model.Ingester {
	for name, c := range types {
		// Find header and if missing (and type is generated) batch insert.
		_, ok := lo.Find(headers, func(h string) bool {
			return h == name
		})

		if !ok && c.IsGenerated {
			return &model.BatchInsertIngester{}
		}
	}

	return &model.CopyFromIngester{}
}

func fileLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanning file: %w", err)
	}

	return lines, nil
}
