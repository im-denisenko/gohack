package main

import (
	"os"
	"testing"
)

func TestResolvePath(t *testing.T) {
	cwd, err := os.Getwd()
	Check(err)

	relativePath := "input/a.json"
	absolutePath := cwd + "/input/a.json"

	relativePathResolved, _ := ResolvePath(relativePath)
	absolutePathResolved, _ := ResolvePath(absolutePath)

	if relativePathResolved != cwd+"/"+relativePath {
		t.Fatalf("relative path resolve was incorrect")
	}
	if absolutePathResolved != absolutePath {
		t.Fatalf("absolute path resolve was incorrect")
	}
}

func TestReportUpdateGroupByCategory(t *testing.T) {
	report := make(Report)

	report.Update(Transaction{UserId: 10, Amount: 10, Category: "aa"})
	report.Update(Transaction{UserId: 10, Amount: 20, Category: "aa"})
	report.Update(Transaction{UserId: 10, Amount: 40, Category: "bb"})

	if len(report) != 1 {
		t.Fatalf("unexpected report length")
	}

	if report[10]["user_id"] != 10 {
		t.Fatalf("user_id was set incorrect")
	}

	if report[10]["sum"] != 70 {
		t.Fatalf("sum was set incorrect")
	}

	if report[10]["category_aa"] != 30 {
		t.Fatalf("category_aa was set incorrect")
	}

	if report[10]["category_bb"] != 40 {
		t.Fatalf("category_bb was set incorrect")
	}
}

func TestReportUpdateGroupByUser(t *testing.T) {
	report := make(Report)

	report.Update(Transaction{UserId: 10, Amount: 10, Category: "aa"})
	report.Update(Transaction{UserId: 10, Amount: 20, Category: "aa"})
	report.Update(Transaction{UserId: 11, Amount: 30, Category: "bb"})
	report.Update(Transaction{UserId: 11, Amount: 40, Category: "bb"})

	if len(report) != 2 {
		t.Fatalf("unexpected report length")
	}

	if report[10]["user_id"] != 10 {
		t.Fatalf("user_id was set incorrect")
	}

	if report[10]["sum"] != 30 {
		t.Fatalf("sum was set incorrect")
	}

	if report[10]["category_aa"] != 30 {
		t.Fatalf("category_aa was set incorrect")
	}

	if report[11]["user_id"] != 11 {
		t.Fatalf("user_id was set incorrect")
	}

	if report[11]["sum"] != 70 {
		t.Fatalf("sum was set incorrect")
	}

	if report[11]["category_bb"] != 70 {
		t.Fatalf("category_aa was set incorrect")
	}
}

func TestReportUpdateGroupEmptyCategory(t *testing.T) {
	report := make(Report)

	report.Update(Transaction{UserId: 10, Amount: 10})

	if len(report) != 1 {
		t.Fatalf("unexpected report length")
	}

	if report[10]["user_id"] != 10 {
		t.Fatalf("user_id was set incorrect")
	}

	if report[10]["sum"] != 10 {
		t.Fatalf("sum was set incorrect")
	}

	if report[10]["category_"] != 10 {
		t.Fatalf("category_ was set incorrect")
	}
}
