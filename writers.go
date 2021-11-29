package main

// -----------------------------------------------------------------------------
// Imports
// -----------------------------------------------------------------------------
import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

// -----------------------------------------------------------------------------
// ReportWriter
// -----------------------------------------------------------------------------
type ReportWriter interface {
	Write(report Report, outputPath string) error
}

// -----------------------------------------------------------------------------
// ReportWriterCSV
// -----------------------------------------------------------------------------
type ReportWriterCSV struct {
}

func (self ReportWriterCSV) Write(report Report, outputPath string) error {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	headersMap := make(map[string]bool)

	for k := range report {
		for kk := range report[k] {
			headersMap[kk] = true
		}
	}

	headers := make([]string, 0, len(headersMap))

	for k := range headersMap {
		headers = append(headers, k)
	}

	w.Write(headers)

	for k := range report {
		values := make([]string, 0, len(report[k]))
		for _, kk := range headers {
			if v, found := report[k][kk]; found {
				values = append(values, fmt.Sprint(v))
			}
		}
		w.Write(values)
	}

	w.Flush()

	return os.WriteFile(outputPath, []byte(b.String()), 0644)
}

// -----------------------------------------------------------------------------
// ReportWriterJSON
// -----------------------------------------------------------------------------
type ReportWriterJSON struct {
}

func (self ReportWriterJSON) Write(report Report, outputPath string) error {
	txs := make([]ReportRow, 0, len(report))

	for _, tx := range report {
		txs = append(txs, tx)
	}

	s, err := json.MarshalIndent(txs, "", "\t")

	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(string(s)), 0644)
}

// -----------------------------------------------------------------------------
// ReportWriterSQLite
// -----------------------------------------------------------------------------
type ReportWriterSQLite struct {
}

func (self ReportWriterSQLite) Write(report Report, outputPath string) error {
	db, err := sql.Open("sqlite3", outputPath)

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
		DROP TABLE IF EXISTS users_report;

		DROP TABLE IF EXISTS users_categories_report;

		CREATE TABLE users_report (
			id INT NOT NULL,
			sum INT NOT NULL,
			PRIMARY KEY (id)
		);

		CREATE TABLE users_categories_report (
			user_id INT NOT NULL,
			category_name VARCHAR(255) NOT NULL,
			sum INT NOT NULL,
			PRIMARY KEY (user_id, category_name)
		);
	`)

	if err != nil {
		return err
	}

	for _, row := range report {
		_, err = db.Exec(`
			INSERT INTO users_report (id, sum)
			VALUES (?, ?)
		`, row["user_id"], row["sum"])

		if err != nil {
			return err
		}

		for kk, vv := range row {
			if strings.HasPrefix(kk, "category_") {
				_, err = db.Exec(`
					INSERT INTO users_categories_report (user_id, category_name, sum)
					VALUES (?, ?, ?)
				`, row["user_id"], strings.TrimPrefix(kk, "category_"), vv)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// CreateReportWriter
// -----------------------------------------------------------------------------
func CreateReportWriter(format string) (ReportWriter, error) {
	if "csv" == format {
		return new(ReportWriterCSV), nil
	}
	if "json" == format {
		return new(ReportWriterJSON), nil
	}
	if "sqlite" == format {
		return new(ReportWriterSQLite), nil
	}
	return nil, errors.New("unknown format: " + format)
}
