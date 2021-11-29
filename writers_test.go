package main

import (
	"testing"
)

func TestCreateReportWriterCSV(t *testing.T) {
	writer, err := CreateReportWriter("csv")

	if err != nil {
		t.Fatal(err)
	}

	_, ok := writer.(*ReportWriterCSV)

	if !ok {
		t.Fatalf("create csv writer was incorrect")
	}
}

func TestCreateReportWriterJSON(t *testing.T) {
	writer, err := CreateReportWriter("json")

	if err != nil {
		t.Fatal(err)
	}

	_, ok := writer.(*ReportWriterJSON)

	if !ok {
		t.Fatalf("create json writer was incorrect")
	}
}

func TestCreateReportWriterSQLite(t *testing.T) {
	writer, err := CreateReportWriter("sqlite")

	if err != nil {
		t.Fatal(err)
	}

	_, ok := writer.(*ReportWriterSQLite)

	if !ok {
		t.Fatalf("create json writer was incorrect")
	}
}

func TestCreateReportWriterUnknown(t *testing.T) {
	_, err := CreateReportWriter("aaaaaa")

	if err == nil {
		t.Errorf("program was not stopped after unknown writer")
	}
}
