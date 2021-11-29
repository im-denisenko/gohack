package main

import (
	"testing"
)

func TestCreateReportGeneratorNaive(t *testing.T) {
	generator, err := CreateReportGenerator("naive")

	if err != nil {
		t.Fatal(err)
	}

	_, ok := generator.(*ReportGeneratorNaive)

	if !ok {
		t.Fatalf("create csv writer was incorrect")
	}
}

func TestCreateReportGeneratorStream(t *testing.T) {
	generator, err := CreateReportGenerator("stream")

	if err != nil {
		t.Fatal(err)
	}

	_, ok := generator.(*ReportGeneratorStream)

	if !ok {
		t.Fatalf("create json writer was incorrect")
	}
}

func TestCreateReportGeneratorUnknown(t *testing.T) {
	_, err := CreateReportGenerator("aaaaaa")

	if err == nil {
		t.Errorf("program was not stopped after unknown writer")
	}
}
