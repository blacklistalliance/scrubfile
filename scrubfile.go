package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	file                  = flag.String("file", "", "File name to scrub (REQUIRED)")
	apiKey                = flag.String("apikey", "", "API Key put between single quotes (REQUIRED) ")
	colnum                = flag.String("colnum", "1", "Column number to use")
	splitchar             = flag.String("splitchar", ",", "What character do we split on wrap in single quotes")
	hasheader             = flag.Bool("hasheader", false, "Does the file have a header")
	download_carrier_data = flag.Bool("carrier", false, "Include carrier data")
	download_invalid      = flag.Bool("invalid", false, "Include invalid data")
	download_no_carrier   = flag.Bool("noCarrier", false, "Include no carrier")
	download_wireless     = flag.Bool("wireless", false, "Include wireless")
	download_federal_dnc  = flag.Bool("federaldnc", false, "Include federal DNC")
	include_feeds         = flag.Bool("includefeeds", false, "Include carrier feeds")
)

const UPLOAD_URL = "https://api.blacklistalliance.net/bulk/upload"

type zipfileinfo struct {
	name string
	size int64
	mod  time.Time
}

/*
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
*/
func (z zipfileinfo) Name() string {
	return z.name
}
func (z zipfileinfo) Size() int64 {
	return z.size
}
func (z zipfileinfo) Mode() fs.FileMode {
	return 0666
}
func (z zipfileinfo) ModTime() time.Time {
	return z.mod
}
func (z zipfileinfo) Sys() interface{} {
	return nil
}
func (z zipfileinfo) IsDir() bool {
	return false
}

func main() {
	fmt.Printf("---- Blacklist Alliance Scrub File Application ----\n")
	flag.Parse()
	if len(*file) == 0 || len(*apiKey) == 0 {
		fmt.Printf("INVALID OPTION PLEASE SEE BELOW\n")
		flag.PrintDefaults()
		return
	}
	params := map[string]string{
		"filetype":              "csv",
		"colnum":                *colnum,
		"splitchar":             *splitchar,
		"key":                   *apiKey,
		"hasheader":             strconv.FormatBool(*hasheader),
		"download_carrier_data": strconv.FormatBool(*download_carrier_data),
		"download_invalid":      strconv.FormatBool(*download_invalid),
		"download_no_carrier":   strconv.FormatBool(*download_no_carrier),
		"download_wireless":     strconv.FormatBool(*download_wireless),
		"download_federal_dnc":  strconv.FormatBool(*download_federal_dnc),
	}
	err := processFile(*file, params, *include_feeds)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	os.Exit(1)
}

func processFile(myfile string, params map[string]string, includefeeds bool) error {
	zipbytes, err := getZipBytes(myfile)
	if err != nil {
		return err
	}

	zipname := strings.Replace(filepath.Base(myfile), filepath.Ext(myfile), ".zip", 1)
	req, err := newfileUploadRequest(UPLOAD_URL, params, zipname, zipbytes)
	if err != nil {
		return fmt.Errorf("Failed to create upload request \n\tError: %v\n", err)
	}
	retBytes, err := makeRequest(req)
	if err != nil {
		return fmt.Errorf("Request failed\n\tError: %v\n", err)
	}
	fmt.Printf("Returned %v bytes\n", len(retBytes))
	err = unzipSaveResponse(retBytes, filepath.Base(myfile), includefeeds)
	if err != nil {
		return fmt.Errorf("Failed to save response\n\tError: %v\n", err)
	}
	return nil
}

func unzipSaveResponse(ret []byte, orgfilename string, includefeeds bool) error {
	retReader := bytes.NewReader(ret)
	zipReader, err := zip.NewReader(retReader, int64(len(ret)))
	if err != nil {
		return err
	}
	oname := filepath.Base(orgfilename)
	oext := filepath.Ext(orgfilename)
	oname = strings.Replace(oname, oext, "", 1)
	if false {
		err = os.WriteFile(oname+".zip", ret, 0666)
		if err != nil {
			return err
		}
	}

	for _, f := range zipReader.File {
		fmt.Printf("Received file %v\n", f.Name)
		if f.Name == "included_feeds.txt" && !includefeeds {
			continue
		}
		justname := strings.Replace(filepath.Base(f.Name), filepath.Ext(f.Name), "", 1)
		nfilename := justname + "_" + f.Name
		if f.Name != "carrier.csv" && f.Name != "included_feeds.txt" {
			nfilename = oname + "_" + justname + oext
		}
		freader, err := f.Open()
		if err != nil {
			return err
		}
		fb, err := ioutil.ReadAll(freader)
		freader.Close()
		err = os.WriteFile(nfilename, fb, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeRequest(req *http.Request) ([]byte, error) {
	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to make http request\n\tError: %v\n", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response\n\tError: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Http request failed not 200\n\tMsg: %v\n", string(respBytes))
	}
	return respBytes, nil
}

func newfileUploadRequest(uri string, params map[string]string, filename string, filebytes []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	fileBytesReader := bytes.NewReader(filebytes)
	_, err = io.Copy(part, fileBytesReader)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func getZipBytes(fileName string) ([]byte, error) {
	fi, err := os.Stat(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %v\n\tError: %v\n", fileName, err)
	}
	if fi.IsDir() {
		fmt.Printf("You must specify a file not a directory\n")
	}
	if fi.Size() == 0 {
		fmt.Printf("Empty file\n")
	}
	fmt.Printf("Checking file %v with a size of %v bytes\n", fi.Name(), fi.Size())
	zipBuffer := bytes.Buffer{}
	zipWriter := zip.NewWriter(&zipBuffer)
	numbersinfo := zipfileinfo{name: fi.Name(), mod: time.Now(), size: int64(unsafe.Sizeof(fi.Size()))}
	zipFileHeader, err := zip.FileInfoHeader(numbersinfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to create zip file\n\tError: %v\n", err)
	}
	zipFileHeader.Name = numbersinfo.name
	zipFileHeader.Method = zip.Deflate
	zipIOWriter, err := zipWriter.CreateHeader(zipFileHeader)
	if err != nil {
		return nil, fmt.Errorf("Unknown zip error\n\tError: %v\n", err)
	}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file\n\tError: %v\n", err)
	}
	_, err = io.Copy(zipIOWriter, f)
	if err != nil {
		return nil, fmt.Errorf("Failed to copy file to zip\n\tError: %v\n", err)
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to write zip file\n\tError: %v\n", err)
	}
	zipbytes := zipBuffer.Bytes()
	fmt.Printf("- Zipped file %v Original Size: %v bytes, New Size: %v bytes\n", fi.Name(), fi.Size(), len(zipbytes))
	if false {
		_ = os.WriteFile("test.zip", zipbytes, 0666)
	}
	return zipbytes, nil
}