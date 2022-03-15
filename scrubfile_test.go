package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestProcessFile(t *testing.T) {
	apikey := os.Getenv("APIKEY")
	if len(apikey) == 0 {
		t.Fatal("Missing APIKEY variable")
	}
	testnum := os.Getenv("TESTNUMBERS")
	testnums := strings.Split(testnum, ",")
	filename := "test.csv"
	filecontents := strings.Join(testnums, "\n")
	err := os.WriteFile(filename, []byte(filecontents), 0666)
	if err != nil {
		t.Fatalf("Unable to write file %v", err)
	}
	params := map[string]string{
		"filetype":              "csv",
		"colnum":                "1",
		"splitchar":             ",",
		"key":                   apikey,
		"hasheader":             strconv.FormatBool(false),
		"download_carrier_data": strconv.FormatBool(false),
		"download_invalid":      strconv.FormatBool(false),
		"download_no_carrier":   strconv.FormatBool(false),
		"download_wireless":     strconv.FormatBool(false),
		"download_federal_dnc":  strconv.FormatBool(false),
	}
	err = processFile(filename, params, *include_feeds)
	if err != nil {
		t.Fatalf("Failed to processfile \n\tError: %v", err)
	}
	fb, err := os.ReadFile("test_all_clean.csv")
	if err != nil {
		t.Fatalf("Failed to read all clean file %v", err)
	}
	fbs := string(fb)
	fbsa := strings.Split(fbs, "\n")
	if len(fbsa) != 2 {
		t.Fatalf("Returned more then one number %v", fbsa)
	}
	if fbsa[0] != "2132133000" {
		t.Fatalf("Wrong number returend %v", fbsa[0])
	}
}