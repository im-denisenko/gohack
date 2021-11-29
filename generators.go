package main

// -----------------------------------------------------------------------------
// Imports
// -----------------------------------------------------------------------------
import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// -----------------------------------------------------------------------------
// ReportGenerator
// -----------------------------------------------------------------------------
type ReportGenerator interface {
	Generate(fp *os.File) (Report, error)
}

// -----------------------------------------------------------------------------
// ReportGeneratorNaive
// -----------------------------------------------------------------------------
type ReportGeneratorNaive struct {
}

func (self ReportGeneratorNaive) Generate(fp *os.File) (Report, error) {
	byteValue, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	Checkpoint("json readed")

	var transactions []Transaction
	err = json.Unmarshal(byteValue, &transactions)
	if err != nil {
		return nil, err
	}
	Checkpoint("json parsed")

	report := make(Report)
	for i := 0; i < len(transactions); i++ {
		report.Update(transactions[i])
	}

	return report, nil

}

// -----------------------------------------------------------------------------
// ReportGeneratorStream
// -----------------------------------------------------------------------------
type ReportGeneratorStream struct {
}

func (self ReportGeneratorStream) Generate(fp *os.File) (Report, error) {
	report := make(Report)
	decoder := json.NewDecoder(fp)

	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return nil, errors.New("expected [")
	}

	for decoder.More() {
		transaction := Transaction{}
		err := decoder.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		report.Update(transaction)
	}

	token, err = decoder.Token()
	if err != nil {
		return nil, err
	}
	if delim, ok := token.(json.Delim); !ok || delim != ']' {
		return nil, errors.New("expected ]")
	}

	return report, nil
}

// -----------------------------------------------------------------------------
// CreateReportGenerator
// -----------------------------------------------------------------------------
func CreateReportGenerator(algorithm string) (ReportGenerator, error) {
	if "naive" == algorithm {
		return new(ReportGeneratorNaive), nil
	}
	if "stream" == algorithm {
		return new(ReportGeneratorStream), nil
	}
	return nil, errors.New("unknown algorithm: " + algorithm)
}
