package main

// -----------------------------------------------------------------------------
// Imports
// -----------------------------------------------------------------------------
import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

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
// Utils
// -----------------------------------------------------------------------------
var start = time.Now()
var helpArg = flag.Bool("h", false, "Print help and exit")
var algoArg = flag.String("a", "naive", "Algorithm to use: naive, stream")
var inputArg = flag.String("i", "input/10M.json", "Path to the input file")
var outputArg = flag.String("o", "output/report.json", "Path to the output file")
var formatArg = flag.String("f", "json", "Report format: json, csv, sqlite")
var quietArg = flag.Bool("q", false, "Disable logs output")

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Checkpoint(message string) {
	if *quietArg == false {
		log.Println(time.Since(start).String() + " " + message)
	}
}

func ResolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, path), nil
}

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------
func main() {
	flag.Parse()
	if (*inputArg == "") || (*outputArg == "") || (*formatArg == "") || (*helpArg == true) {
		flag.PrintDefaults()

		if *helpArg == true {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
	Checkpoint("program started")

	inputPath, err := ResolvePath(*inputArg)
	Check(err)

	inputFile, err := os.Open(inputPath)
	Check(err)

	defer inputFile.Close()
	Checkpoint("input opened")

	generator, err := CreateReportGenerator(*algoArg)
	Check(err)

	report, err := generator.Generate(inputFile)
	Check(err)
	Checkpoint("report generated")

	writer, err := CreateReportWriter(*formatArg)
	Check(err)

	outputPath, err := ResolvePath(*outputArg)
	Check(err)

	Check(os.MkdirAll(filepath.Dir(outputPath), 0755))
	Check(writer.Write(report, outputPath))
	Checkpoint("report persisted")

	Checkpoint("program finished")
}
