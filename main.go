package main

// -----------------------------------------------------------------------------
// Imports
// -----------------------------------------------------------------------------
import "bytes"
import "encoding/csv"
import "encoding/json"
import "errors"
import "flag"
import "fmt"
import "log"
import "os"
import "path/filepath"
import "time"
import "io/ioutil"

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------
type Report map[int]ReportRow

type ReportRow map[string]int

type Transaction struct {
	UserId   int    `json:"user_id"`
	Amount   int    `json:"amount"`
	Category string `json:"category"`
}

// -----------------------------------------------------------------------------
// ReportGenerator
// -----------------------------------------------------------------------------
type ReportGenerator interface {
	Generate(fp *os.File) Report
}

// -----------------------------------------------------------------------------
// ReportGeneratorNaive
// -----------------------------------------------------------------------------
type ReportGeneratorNaive struct {
}

func (self ReportGeneratorNaive) Generate(fp *os.File) Report {
	byteValue, _ := ioutil.ReadAll(fp)
	Checkpoint("json readed")

	var transactions []Transaction
	json.Unmarshal(byteValue, &transactions)
	Checkpoint("json parsed")

	report := make(Report)
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		if report[transaction.UserId] == nil {
			report[transaction.UserId] = make(ReportRow)
		}
		report[transaction.UserId]["sum"] += transaction.Amount
		report[transaction.UserId]["user_id"] = transaction.UserId
		report[transaction.UserId]["category_"+transaction.Category] += transaction.Amount
	}

	return report

}

// -----------------------------------------------------------------------------
// ReportGeneratorStream
// -----------------------------------------------------------------------------
type ReportGeneratorStream struct {
}

func (self ReportGeneratorStream) Generate(fp *os.File) Report {
	report := make(Report)
	decoder := json.NewDecoder(fp)

	token, err := decoder.Token()
	Check(err)
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		Check(errors.New("expected ["))
	}

	for decoder.More() {
		transaction := Transaction{}
		Check(decoder.Decode(&transaction))

		if report[transaction.UserId] == nil {
			report[transaction.UserId] = make(ReportRow)
		}

		report[transaction.UserId]["sum"] += transaction.Amount
		report[transaction.UserId]["user_id"] = transaction.UserId
		report[transaction.UserId]["category_"+transaction.Category] += transaction.Amount
	}

	token, err = decoder.Token()
	Check(err)
	if delim, ok := token.(json.Delim); !ok || delim != ']' {
		Check(errors.New("expected ]"))
	}

	return report
}

// -----------------------------------------------------------------------------
// CreateReportGenerator
// -----------------------------------------------------------------------------
func CreateReportGenerator(algorithm string) ReportGenerator {
	if "naive" == algorithm {
		return new(ReportGeneratorNaive)
	}
	if "stream" == algorithm {
		return new(ReportGeneratorStream)
	}
	Check(errors.New("unknown algorithm: " + algorithm))
	return nil
}

// -----------------------------------------------------------------------------
// ReportFormatter
// -----------------------------------------------------------------------------
type ReportFormatter interface {
	Format(report Report) string
}

// -----------------------------------------------------------------------------
// ReportFormatterCSV
// -----------------------------------------------------------------------------
type ReportFormatterCSV struct {
}

func (self ReportFormatterCSV) Format(report Report) string {
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

	return b.String()
}

// -----------------------------------------------------------------------------
// ReportFormatterJSON
// -----------------------------------------------------------------------------
type ReportFormatterJSON struct {
}

func (self ReportFormatterJSON) Format(report Report) string {
	txs := make([]ReportRow, 0, len(report))
	for _, tx := range report {
		txs = append(txs, tx)
	}
	s, err := json.MarshalIndent(txs, "", "\t")
	Check(err)
	return string(s)
}

// -----------------------------------------------------------------------------
// CreateReportFormatter
// -----------------------------------------------------------------------------
func CreateReportFormatter(format string) ReportFormatter {
	if "csv" == format {
		return new(ReportFormatterCSV)
	}
	if "json" == format {
		return new(ReportFormatterJSON)
	}
	Check(errors.New("unknown format: " + format))
	return nil
}

// -----------------------------------------------------------------------------
// Utils
// -----------------------------------------------------------------------------
var start = time.Now()

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Checkpoint(message string) {
	log.Println(time.Since(start).String() + " " + message)
}

func ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	cwd, err := os.Getwd()
	Check(err)
	return filepath.Join(cwd, path)
}

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------
func main() {
	Checkpoint("program started")

	helpArg := flag.Bool("h", false, "Print help and exit")
	algoArg := flag.String("a", "naive", "Algorithm to use: naive, stream")
	inputArg := flag.String("i", "input/10M.json", "Path to input file")
	outputArg := flag.String("o", "output/report.json", "Path to output file")
	formatArg := flag.String("f", "json", "Report format: json, csv")
	flag.Parse()
	Checkpoint("flags parsed")

	if (*inputArg == "") || (*outputArg == "") || (*formatArg == "") || (*helpArg == true) {
		flag.PrintDefaults()
		os.Exit(1)
	}
	Checkpoint("flags validated")

	inputPath := ResolvePath(*inputArg)
	inputFile, err := os.Open(inputPath)
	Check(err)
	defer inputFile.Close()
	Checkpoint("input opened")

	generator := CreateReportGenerator(*algoArg)
	report := generator.Generate(inputFile)
	Checkpoint("report generated")

	formatter := CreateReportFormatter(*formatArg)
	formattedReport := formatter.Format(report)
	Checkpoint("report formatted")

	outputPath := ResolvePath(*outputArg)
	Check(os.MkdirAll(filepath.Dir(outputPath), 0755))
	Check(os.WriteFile(outputPath, []byte(formattedReport), 0644))
	Checkpoint("program finished")
}
